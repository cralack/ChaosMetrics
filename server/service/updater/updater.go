package updater

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/pkg/fetcher"
	"github.com/cralack/ChaosMetrics/server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Updater interface {
	Update(loc uint, args ...string) error
}

type RiotUpdater struct {
	Logger *zap.Logger
	Lock   *sync.Mutex
	DB     *gorm.DB
	
	MatchVis map[string]struct{}
	Fetcher  fetcher.Fetcher
	// Conf     UpdateConfig
}

func NewRiotUpdater() *RiotUpdater {
	return &RiotUpdater{
		Logger:   global.GVA_LOG,
		Lock:     &sync.Mutex{},
		DB:       global.GVA_DB,
		MatchVis: make(map[string]struct{}),
		Fetcher:  fetcher.NewBrowserFetcher(),
	}
}

func (u *RiotUpdater) UpdateSummonerMatch(loc uint, puuid string, cnt int) (res []*riotmodel.MatchDto, err error) {
	// get match IDs
	matchIDs, err := u.getMatchIDsByPUUID(loc, puuid, cnt)
	if err != nil {
		return nil, err
	}
	matches := make([]*riotmodel.MatchDto, 0, cnt)
	for _, id := range matchIDs {
		mat, err := u.getMatchByMatchID(loc, id)
		if err != nil {
			return nil, err
		}
		if mat != nil && err == nil {
			matches = append(matches, mat)
		}
	}
	
	return matches, nil
}

// Get a list of match ids by puuid
// api: /lol/match/v5/matches/by-puuid/{puuid}/ids
func (u *RiotUpdater) getMatchIDsByPUUID(loc uint, puuid string, count int) (res []string, err error) {
	// set up args val
	prefix := utils.ConvertPlatformToHost(loc)
	stem := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids?", puuid)
	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
	endTime := time.Now().Unix()                     // cur time unix
	suffix := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d",
		startTime, endTime, count)
	// fill buffer
	u.Lock.Lock()
	defer u.Lock.Unlock()
	
	buff, err := u.Fetcher.Get(prefix + stem + suffix)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("fetch %s failed", stem),
			zap.Error(err))
		return nil, err
	}
	// parse
	err = json.Unmarshal(buff, &res)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
			"matchIDs"), zap.Error(err))
		return nil, err
	}
	u.Logger.Info("get matchIDs by puuid succeed")
	return res, nil
}

// Get a match by match id
// api: /lol/match/v5/matches/{matchId}
func (u *RiotUpdater) getMatchByMatchID(loc uint, matchID string) (res *riotmodel.MatchDto, err error) {
	if _, has := u.MatchVis[matchID]; has {
		return nil, nil
	}
	// set up args val
	prefix := utils.ConvertPlatformToHost(loc)
	stem := fmt.Sprintf("/lol/match/v5/matches/%s", matchID)
	// fetch buffer
	u.Lock.Lock()
	defer u.Lock.Unlock()
	
	buff, err := u.Fetcher.Get(prefix + stem)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("fetch %s failed", stem),
			zap.Error(err))
		return nil, err
	}
	// parse
	if err = json.Unmarshal(buff, &res); err != nil {
		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
			"MatchDTO"), zap.Error(err))
		return nil, err
	}
	// store
	u.DB.Save(res)
	return res, nil
	
}

// Get the BEST(challenger/grandmaster/master) league for given queue
// api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
func (u *RiotUpdater) UpdateBetserLeague(loc, tier, que uint) (res []*riotmodel.LeagueEntryDTO, err error) {
	var stem string
	var tierStr string
	prefix := utils.ConvertPlatformURL(loc)
	queStr := getQueueString(que)
	switch tier {
	case riotmodel.CHALLENGER:
		tierStr = "CHALLENGER"
	case riotmodel.GRANDMASTER:
		tierStr = "GRANDMASTER"
	case riotmodel.MASTER:
		tierStr = "MASTER"
	}
	stem = fmt.Sprintf("/lol/league/v4/%sleagues/by-queue/%s",
		strings.ToLower(tierStr), queStr)
	// fill buffer
	u.Lock.Lock()
	defer u.Lock.Unlock()
	buff, err := u.Fetcher.Get(prefix + stem)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("fetch %s failed", stem),
			zap.Error(err))
		return nil, err
	}
	// parse
	results := &riotmodel.LeagueListDTO{}
	if err = json.Unmarshal(buff, &results); err != nil {
		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
			"LeagueListDTO"), zap.Error(err))
		return nil, err
	}
	res = results.Entries
	for _, enrty := range res {
		enrty.Tier = tierStr
	}
	// todo puller
	u.DB.Save(res)
	return res, nil
}

// Get all the league entries
// api: /lol/league/v4/entries/{queue}/{tier}/{division}
// page(optional): 	Defaults to 1. Starts with page 1.
func (u *RiotUpdater) UpdateMortalLeague(loc, tier, division, que, idx uint) (res []*riotmodel.LeagueEntryDTO, err error) {
	prefix := utils.ConvertPlatformURL(loc)
	queStr := getQueueString(que)
	tierDiv := getMortalString(tier, division)
	//todo:page should be int
	// page := strconv.Itoa(int(idx))
	// stem := fmt.Sprintf("/lol/league/v4/entries/%s/%s/%s?page=%d",
	// 	queStr, tierDiv[0], tierDiv[1], page)
	stem := fmt.Sprintf("/lol/league/v4/entries/%s/%s/%s",
		queStr, tierDiv[0], tierDiv[1])
	// fill buffer
	u.Lock.Lock()
	defer u.Lock.Unlock()
	buff, err := u.Fetcher.Get(prefix + stem)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("fetch %s failed", stem),
			zap.Error(err))
		return nil, err
	}
	// parse
	if err = json.Unmarshal(buff, &res); err != nil {
		u.Logger.Error(fmt.Sprintf("unmarshal json to %s failed",
			"LeagueItemDTO"), zap.Error(err))
		return nil, err
	}
	for _, enrty := range res {
		enrty.Tier = tierDiv[0]
	}
	// todo puller
	u.DB.Save(res)
	return res, nil
}

func getQueueString(que uint) string {
	switch que {
	case riotmodel.RANKED_SOLO_5x5:
		return "RANKED_SOLO_5x5"
	case riotmodel.RANKED_FLEX_SR:
		return "RANKED_FLEX_SR"
	case riotmodel.RANKED_FLEX_TT:
		return "RANKED_FLEX_TT"
	default:
		return ""
	}
}

func getMortalString(tier, div uint) (ans []string) {
	ans = make([]string, 0, 2)
	switch tier {
	case riotmodel.DIAMOND:
		ans = append(ans, "DIAMOND")
	case riotmodel.PLATINUM:
		ans = append(ans, "PLATINUM")
	case riotmodel.GOLD:
		ans = append(ans, "GOLD")
	case riotmodel.SILVER:
		ans = append(ans, "SILVER")
	case riotmodel.BRONZE:
		ans = append(ans, "BRONZE")
	case riotmodel.IRON:
		ans = append(ans, "IRON")
	default:
		return
	}
	switch div {
	case 1:
		ans = append(ans, "I")
	case 2:
		ans = append(ans, "II")
	case 3:
		ans = append(ans, "III")
	case 4:
		ans = append(ans, "IV")
	default:
		return
	}
	return
}

type Task struct {
	key    string
	brief  string
	buffer interface{}
}

func (u *RiotUpdater) Puller(req chan Task) {
	for {
		select {
		case task := <-req:
			u.Logger.Info(fmt.Sprintf("pulling task type:%s,brief:%s\n",
				task.key, task.brief))
			u.Syncer(task.key, task.buffer)
			u.Logger.Info(fmt.Sprintf("task type:%s,brief:%s store succeed\n",
				task.key, task.brief))
		}
	}
}

func (u *RiotUpdater) Syncer(key string, data interface{}) {
	switch key {
	// case "summoner":
	// 	summoners, ok := data.([]*riotmodel.SummonerDTO)
	// 	if !ok {
	// 		u.Logger.Error("buffer'key and buff doens match")
	// 	}
	// 	u.DB.Save(summoners)
	case "match":
		matches, ok := data.([]*riotmodel.MatchDto)
		if !ok {
			u.Logger.Error("buffer'key and buff doens match")
		}
		u.DB.Save(matches)
	case "leagueEntry":
		entries, ok := data.([]*riotmodel.LeagueEntryDTO)
		if !ok {
			u.Logger.Error("buffer'key and buff doens match")
		}
		u.DB.Save(entries)
	}
}

// func (u *RiotUpdater) matchesStoreer(data interface{}) error {
// 	matches, ok := data.([]*riotmodel.MatchDto)
//
// }
