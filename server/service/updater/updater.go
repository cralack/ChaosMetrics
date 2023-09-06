package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// type Updater interface {
// 	Update(loc uint, args ...string) error
// }

type Updater struct {
	logger  *zap.Logger
	db      *gorm.DB
	rdb     *redis.Client
	lock    *sync.Mutex
	fetcher fetcher.Fetcher

	CurVersion string
	matchVis   map[string]map[string]bool
	stgy       *Strategy
	// Conf     UpdateConfig
}

func NewRiotUpdater(opts ...Option) *Updater {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}

	return &Updater{
		logger: global.GVA_LOG,
		db:     global.GVA_DB,
		rdb:    global.GVA_RDB,
		lock:   &sync.Mutex{},
		fetcher: fetcher.NewBrowserFetcher(
			fetcher.WithRateLimiter(false),
		),
		matchVis: make(map[string]map[string]bool),
		stgy:     stgy,
	}
}

func (u *Updater) UpdateVersions() (version []string) {
	var (
		buff []byte
		url  string
		err  error
	)
	url = "https://ddragon.leagueoflegends.com/api/versions.json"
	if buff, err = u.fetcher.Get(url); err != nil || buff == nil {
		u.logger.Error("update version failed", zap.Error(err))
	}
	if err = json.Unmarshal(buff, &version); err != nil {
		u.logger.Error("unmarshal json to version failed", zap.Error(err))
	}
	if u.CurVersion == "" {
		u.CurVersion = version[0]
	}
	return
}

// UpdateChampions from Version X.Y.Z
func (u *Updater) UpdateChampions(version string) {
	var (
		buff     []byte
		url      string
		err      error
		flag     bool
		vIdx     uint
		cham     *riotmodel.ChampionDTO
		chamList *riotmodel.ChampionListDTO
	)

	if vIdx, err = utils.ConvertVersionToIdx(version); version == "" || err != nil {
		u.logger.Error("wrong version", zap.Error(err))
	}

	// get champion chamList
	url = fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", version)
	if buff, err = u.fetcher.Get(url); err != nil || buff == nil {
		flag = true
		u.logger.Error("get champion list failed", zap.Error(err))
	}
	if err = json.Unmarshal(buff, &chamList); err != nil {
		flag = true
		u.logger.Error("unmarshal json to champion list failed", zap.Error(err))
	}
	// record each version's champion list
	ctx := context.Background()
	cList := make([]string, 0, len(chamList.Data))
	for cId := range chamList.Data {
		cList = append(cList, cId)
	}
	sort.Strings(cList)
	buff, err = json.Marshal(cList)
	if err != nil {
		flag = true
		u.logger.Error("marshal champion list failed", zap.Error(err))
	}

	if err = u.rdb.HSet(ctx, "/championlist", vIdx, buff).Err(); err != nil {
		flag = true
		u.logger.Error("failed", zap.Error(err))
	}
	// update champion for each lang

	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLanguageCode(langCode)
		key := fmt.Sprintf("/champions/%s", lang)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		cmds := make([]*redis.IntCmd, 0, len(chamList.Data))
		pipe := u.rdb.Pipeline()
		cnt := 0
		for chamID := range chamList.Data {
			// fetch buffer
			url = fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json",
				version, lang, chamID)
			if buff, err = u.fetcher.Get(url); err != nil || buff == nil {
				flag = true
				u.logger.Error(fmt.Sprintf("update %s@%s failed",
					chamID, lang), zap.Error(err))
				continue
			}
			// unmarshal to champion
			var tmp *riotmodel.ChampionSingleDTO
			if err = json.Unmarshal(buff, &tmp); err != nil {
				flag = true
				u.logger.Error(fmt.Sprintf("unmarshal json to %s@%s's model failed",
					chamID, lang), zap.Error(err))
				continue
			} else {
				cnt++
				u.logger.Debug(fmt.Sprintf("fetch %03d/%03d %s@%s succeed",
					cnt, len(chamList.Data), chamID, lang))
			}
			cham = tmp.Data[chamID]
			cmds = append(cmds, pipe.HSet(ctx, key, fmt.Sprintf("%s@%d", cham.ID, vIdx), cham))

		}
		if _, err := pipe.Exec(ctx); err != nil {
			flag = true
			u.logger.Error("redis store champions failed", zap.Error(err))
		}
	}
	if !flag {
		u.logger.Info("all " + version + " champion update done")
	}
	return
}

func (u *Updater) UpdateItems(version string) {

	var (
		buff     []byte
		url      string
		err      error
		flag     bool
		vIdx     uint
		itemList *riotmodel.ItemList
	)
	if version == "" {
		u.logger.Error("wrong version")
	}
	vIdx, err = utils.ConvertVersionToIdx(version)
	ctx := context.Background()
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLanguageCode(langCode)
		key := fmt.Sprintf("/items/%s", lang)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		url = fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/item.json",
			version, lang)
		if buff, err = u.fetcher.Get(url); err != nil {
			flag = true
			u.logger.Error("get item list failed", zap.Error(err))
			continue
		}
		if err = json.Unmarshal(buff, &itemList); err != nil {
			flag = true
			u.logger.Error("unmarshal json to item list failed", zap.Error(err))
			continue
		}
		pipe := u.rdb.Pipeline()
		cmds := make([]*redis.IntCmd, 0, len(itemList.Data))
		for id, item := range itemList.Data {
			item.ID = id
			cmds = append(cmds, pipe.HSet(ctx, key, fmt.Sprintf("%s@%d", id, vIdx), item))
		}

		if _, err := pipe.Exec(ctx); err != nil {
			flag = true
			u.logger.Error("redis store items failed", zap.Error(err))
		}
	}
	if !flag {
		u.logger.Info("all " + version + "'s items update done")
	}
	return
}

func (u *Updater) UpdatePerks() {
	var (
		buff        []byte
		url         string
		err         error
		key         string
		perks       []*riotmodel.Perk
		perkDetails []*riotmodel.PerkDetail
	)
	ctx := context.Background()
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLanguageCode(langCode)
		key = fmt.Sprintf("/perks/%s", lang)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		// fetch perk relation
		url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/runesReforged.json",
			u.CurVersion, lang)
		if buff, err = u.fetcher.Get(url); err != nil {
			u.logger.Error("get  perks failed", zap.Error(err))
			continue
		}
		if err = json.Unmarshal(buff, &perks); err != nil {
			u.logger.Error("unmarshal json to item list failed", zap.Error(err))
			continue
		}
		// store perk relation
		for _, p := range perks {
			if err = u.rdb.HSet(ctx, key, p.ID, p).Err(); err != nil {
				u.logger.Error("save perks failed", zap.Error(err))
				break
			}
		}
		// fetch perk detail
		if lang == "en_US" {
			lang = "default"
		}
		url = fmt.Sprintf("https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/%s/v1/perks.json",
			strings.ToLower(lang))
		if buff, err = u.fetcher.Get(url); err != nil {
			u.logger.Error("get  perks failed", zap.Error(err))
			continue
		}
		if err = json.Unmarshal(buff, &perkDetails); err != nil {
			u.logger.Error("unmarshal json to item list failed", zap.Error(err))
			continue
		}
		for _, p := range perkDetails {
			if err = u.rdb.HSet(ctx, key, p.ID, p).Err(); err != nil {
				u.logger.Error("save perks failed", zap.Error(err))
				break
			}
		}
	}
	u.logger.Info("all perk update done")
}

//
// func (u *Updater) UpdateSummonerMatch(loc uint, puuid string, cnt int) (res []*riotmodel.MatchDTO, err error) {
// 	// get match IDs
// 	matchIDs, err := u.getMatchIDsByPUUID(loc, puuid, cnt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	matches := make([]*riotmodel.MatchDTO, 0, cnt)
// 	for _, id := range matchIDs {
// 		mat, err := u.getMatchByMatchID(loc, id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if mat != nil && err == nil {
// 			matches = append(matches, mat)
// 		}
// 	}
//
// 	return matches, nil
// }
//
// // Get a list of match ids by puuid
// // api: /lol/match/v5/matches/by-puuid/{puuid}/ids
// func (u *Updater) getMatchIDsByPUUID(loc uint, puuid string, count int) (res []string, err error) {
// 	// set up args val
// 	host := utils.ConvertPlatformToHost(loc)
// 	path := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids?", puuid)
// 	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
// 	endTime := time.Now().Unix()                     // cur time unix
// 	suffix := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d",
// 		startTime, endTime, count)
// 	// fill buffer
// 	u.Lock.Lock()
// 	defer u.Lock.Unlock()
//
// 	buff, err := u.Fetcher.Get(host + path + suffix)
// 	if err != nil {
// 		u.Logger.Error(fmt.Sprintf("fetch %s failed", path),
// 			zap.Error(err))
// 		return nil, err
// 	}
// 	// parse
// 	err = json.Unmarshal(buff, &res)
// 	if err != nil {
// 		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
// 			"matchIDs"), zap.Error(err))
// 		return nil, err
// 	}
// 	u.Logger.Info("get matchIDs by puuid succeed")
// 	return res, nil
// }
//
// // Get a match by match id
// // api: /lol/match/v5/matches/{matchId}
// func (u *Updater) getMatchByMatchID(loc uint, matchID string) (res *riotmodel.MatchDTO, err error) {
// 	if _, has := u.MatchVis[matchID]; has {
// 		return nil, nil
// 	}
// 	// set up args val
// 	host := utils.ConvertPlatformToHost(loc)
// 	path := fmt.Sprintf("/lol/match/v5/matches/%s", matchID)
// 	// fetch buffer
// 	u.Lock.Lock()
// 	defer u.Lock.Unlock()
//
// 	buff, err := u.Fetcher.Get(host + path)
// 	if err != nil {
// 		u.Logger.Error(fmt.Sprintf("fetch %s failed", path),
// 			zap.Error(err))
// 		return nil, err
// 	}
// 	// parse
// 	if err = json.Unmarshal(buff, &res); err != nil {
// 		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
// 			"MatchDTO"), zap.Error(err))
// 		return nil, err
// 	}
// 	u.DB.Save(res)
// 	return res, nil
//
// }
//
// // Get the BEST(challenger/grandmaster/master) league for given queue
// // api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
// func (u *Updater) UpdateBetserLeague(loc, tier, que uint) (res []*riotmodel.LeagueEntryDTO, err error) {
// 	var stem string
// 	var tierStr string
// 	_, host := utils.ConvertHostURL(loc)
// 	queStr := getQueueString(que)
// 	switch tier {
// 	case riotmodel.CHALLENGER:
// 		tierStr = "CHALLENGER"
// 	case riotmodel.GRANDMASTER:
// 		tierStr = "GRANDMASTER"
// 	case riotmodel.MASTER:
// 		tierStr = "MASTER"
// 	}
// 	stem = fmt.Sprintf("/lol/league/v4/%sleagues/by-queue/%s",
// 		strings.ToLower(tierStr), queStr)
// 	// fill buffer
// 	u.Lock.Lock()
// 	defer u.Lock.Unlock()
// 	buff, err := u.Fetcher.Get(host + stem)
// 	if err != nil {
// 		u.Logger.Error(fmt.Sprintf("fetch %s failed", stem),
// 			zap.Error(err))
// 		return nil, err
// 	}
// 	// parse
// 	results := &riotmodel.LeagueListDTO{}
// 	if err = json.Unmarshal(buff, &results); err != nil {
// 		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
// 			"LeagueListDTO"), zap.Error(err))
// 		return nil, err
// 	}
// 	res = results.Entries
// 	for _, enrty := range res {
// 		enrty.Tier = tierStr
// 	}
// 	u.DB.Save(res)
// 	return res, nil
// }
//
// // Get all the league entries
// // api: /lol/league/v4/entries/{queue}/{tier}/{division}
// // page(optional): 	Defaults to 1. Starts with page 1.
// func (u *Updater) UpdateMortalLeague(loc, tier, division, que uint) (res []*riotmodel.LeagueEntryDTO, err error) {
// 	_, host := utils.ConvertHostURL(loc)
// 	queStr := getQueueString(que)
// 	rank := getMortalString(tier, division)
// 	page := 0
// 	cnt := 0
// 	u.Lock.Lock()
// 	defer u.Lock.Unlock()
// 	for {
// 		page++
// 		path := fmt.Sprintf("/lol/league/v4/entries/%s/%s/%s?page=%s",
// 			queStr, rank[0], rank[1], strconv.Itoa(page))
// 		// fill buffer
// 		buff, err := u.Fetcher.Get(host + path)
// 		if err != nil {
// 			u.Logger.Error(fmt.Sprintf("fetch %s failed", path),
// 				zap.Error(err))
// 			return nil, err
// 		}
// 		// parse
// 		if err = json.Unmarshal(buff, &res); err != nil {
// 			u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
// 				"LeagueItemDTO"), zap.Error(err))
// 			return nil, err
// 		}
// 		key := fmt.Sprintf("%s_%s", rank[0], rank[1])
// 		if len(res) == 0 {
// 			u.Logger.Info(fmt.Sprintf("all %s data fetch done at page %0d", key, page))
// 			return nil, nil
// 		}
// 		for _, enrty := range res {
// 			enrty.Tier = rank[0]
// 		}
// 		cnt += len(res)
//
// 		u.Schduler.requestCh <- &scheduler.Task{
// 			// Type:    key,
// 			// Brief:  key + ":" + strconv.Itoa(page),
// 			// Buffer: res,
// 		}
// 	}
// 	return nil, nil
// }
//
// func getQueueString(que uint) string {
// 	switch que {
// 	case riotmodel.RANKED_SOLO_5x5:
// 		return "RANKED_SOLO_5x5"
// 	case riotmodel.RANKED_FLEX_SR:
// 		return "RANKED_FLEX_SR"
// 	case riotmodel.RANKED_FLEX_TT:
// 		return "RANKED_FLEX_TT"
// 	default:
// 		return ""
// 	}
// }
//
// func getMortalString(tier, div uint) (ans []string) {
// 	ans = make([]string, 0, 2)
// 	switch tier {
// 	case riotmodel.DIAMOND:
// 		ans = append(ans, "DIAMOND")
// 	case riotmodel.PLATINUM:
// 		ans = append(ans, "PLATINUM")
// 	case riotmodel.GOLD:
// 		ans = append(ans, "GOLD")
// 	case riotmodel.SILVER:
// 		ans = append(ans, "SILVER")
// 	case riotmodel.BRONZE:
// 		ans = append(ans, "BRONZE")
// 	case riotmodel.IRON:
// 		ans = append(ans, "IRON")
// 	default:
// 		return
// 	}
// 	switch div {
// 	case 1:
// 		ans = append(ans, "I")
// 	case 2:
// 		ans = append(ans, "II")
// 	case 3:
// 		ans = append(ans, "III")
// 	case 4:
// 		ans = append(ans, "IV")
// 	default:
// 		return
// 	}
// 	return
// }
//
// func (u *Updater) Syncer(key string, data interface{}) {
// 	switch key {
// 	// case "summoner":
//
// 	case "match":
// 		matches, ok := data.([]*riotmodel.MatchDTO)
// 		if !ok {
// 			u.Logger.Error("buffer'key and buff doens match")
// 		}
// 		u.DB.Save(matches)
// 	case "leagueEntry":
// 		entries, ok := data.([]*riotmodel.LeagueEntryDTO)
// 		if !ok {
// 			u.Logger.Error("buffer'key and buff doens match")
// 		}
// 		u.DB.Save(entries)
// 	}
// }
//
// // func (u *Updater) matchesStoreer(data interface{}) error {
// // 	matches, ok := data.([]*riotmodel.MatchDTO)
// //
// // }
