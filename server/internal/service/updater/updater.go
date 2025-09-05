package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/fetcher"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Updater struct {
	logger  *zap.Logger
	db      *gorm.DB
	rdb     *redis.Client
	lock    *sync.Mutex
	fetcher fetcher.Fetcher

	CurVersion string
	matchVis   map[string]map[string]bool
	stgy       *Strategy
}

func NewRiotUpdater(opts ...Setup) *Updater {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}
	if vers := utils.GetCurMajorVersions(); len(vers) > 0 && stgy.EndMark == "" {
		stgy.EndMark = vers[len(vers)-1]
	}

	return &Updater{
		logger: global.ChaLogger,
		db:     global.ChaDB,
		rdb:    global.ChaRDB,
		lock:   &sync.Mutex{},
		fetcher: fetcher.NewBrowserFetcher(
			fetcher.WithRateLimiter(false),
		),
		matchVis: make(map[string]map[string]bool),
		stgy:     stgy,
	}
}

func (u *Updater) UpdateAll() {
	versions := u.UpdateVersions()
	u.logger.Info("current version:" + versions[0])
	for _, curVer := range versions {
		if isEnd(curVer, u.stgy.EndMark) {
			u.logger.Info(curVer)
			u.UpdatePerks(curVer)
			u.UpdateItems(curVer)
			u.UpdateSpell(curVer)
			u.UpdateChampions(curVer)
		} else {
			u.rdb.HSet(context.Background(), "/lastupdate", "updater", time.Now().Unix())
			break
		}
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

	for i, v := range version {
		if v == "3.6.14" {
			version = version[:i+1]
			break
		}
	}
	// store to redis
	key := "/version"
	buff, _ = json.Marshal(version)
	if err = u.rdb.HSet(context.Background(), key, "versions", buff).Err(); err != nil {
		u.logger.Error("update total version to redis failed", zap.Error(err))
	}
	if err = u.rdb.HSet(context.Background(), key, "cur", version[0]).Err(); err != nil {
		u.logger.Error("update cur version to redis failed", zap.Error(err))
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
		buffer   string
		url      string
		err      error
		vIdx     uint
		cList    []string
		cham     *riotmodel.ChampionDTO
		chamList *riotmodel.ChampionListDTO
	)

	if vIdx, err = utils.ConvertVersionToIdx(version); version == "" || err != nil {
		u.logger.Error("wrong version", zap.Error(err))
	}

	ctx := context.Background()
	cList = make([]string, 0, 300)

	// store idxap[idx]->champion name
	idxMap := make(map[string]string)
	if version == u.CurVersion {
		url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", version)
		if buff, err = u.fetcher.Get(url); err != nil || buff == nil {
			u.logger.Error("get champion list failed", zap.Error(err))
		}
		if err = json.Unmarshal(buff, &chamList); err != nil {
			u.logger.Error("unmarshal json to champion list failed", zap.Error(err))
		}
		for name, c := range chamList.Data {
			idxMap[c.Key] = name
		}
		buff, _ = json.Marshal(idxMap)
		if err = u.rdb.HSet(ctx, "/championlist", "idxmap", buff).Err(); err != nil {
			u.logger.Error("store champion idxmap failed", zap.Error(err))
		}

	}

	// get chamlist from rdb or riot
	if buffer = u.rdb.HGet(ctx, "/championlist", fmt.Sprintf("%d",
		vIdx)).Val(); u.stgy.ForceUpdate || buffer == "" {
		// get champion chamList
		url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", version)
		if buff, err = u.fetcher.Get(url); err != nil || buff == nil {
			u.logger.Error("get champion list failed", zap.Error(err))
		}
		if err = json.Unmarshal(buff, &chamList); err != nil {
			u.logger.Error("unmarshal json to champion list failed", zap.Error(err))
		}
		// record each version's champion list
		for cId := range chamList.Data {
			cList = append(cList, cId)
		}
		// store chamlist
		sort.Strings(cList)
	} else {
		err = json.Unmarshal([]byte(buffer), &cList)
		u.logger.Error("rdb data wrong", zap.Error(err))
	}

	// store champion list
	if buff, err = json.Marshal(cList); err != nil {
		u.logger.Error("marshal champion list failed", zap.Error(err))
	}
	if err = u.rdb.HSet(ctx, "/championlist", fmt.Sprintf("%d", vIdx), buff).Err(); err != nil {
		u.logger.Error("failed", zap.Error(err))
	}

	// update champion for each lang
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLangToLangStr(langCode)
		key := fmt.Sprintf("/champions/%s", lang)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		cmds := make([]*redis.IntCmd, 0, len(cList))
		pipe := u.rdb.Pipeline()
		cnt := 0

		for _, chamID := range cList {
			if buffer = u.rdb.HGet(ctx, key, fmt.Sprintf("%s@%d", chamID,
				vIdx)).Val(); buffer != "" && !u.stgy.ForceUpdate {
				continue
			}

			// fetch buffer
			url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json",
				version, lang, chamID)
			if buff, err = u.fetcher.Get(url); err != nil || buff == nil {
				u.logger.Error(fmt.Sprintf("update %s@%s failed",
					chamID, lang), zap.Error(err))
				continue
			}
			// unmarshal to champion
			var tmp *riotmodel.ChampionSingleDTO
			if err = json.Unmarshal(buff, &tmp); err != nil {
				u.logger.Error(fmt.Sprintf("unmarshal json to %s@%s's model failed",
					chamID, lang), zap.Error(err))
				continue
			} else {
				cnt++
				u.logger.Debug(fmt.Sprintf("fetch %03d/%03d %s@%s succeed",
					cnt, len(cList), chamID, lang))
			}
			cham = tmp.Data[chamID]
			cmds = append(cmds, pipe.HSet(ctx, key, fmt.Sprintf("%s@%d", cham.ID, vIdx), cham))
		}

		// store chamlist
		if _, err = pipe.Exec(ctx); err != nil {
			u.logger.Error("redis store champions failed", zap.Error(err))
		} else {
			u.logger.Info(fmt.Sprintf("%d@%s's champion(%d/%d) store done", vIdx, lang, cnt, len(cList)))
		}
	}
}

func (u *Updater) UpdateItems(version string) {

	var (
		buff     []byte
		url      string
		err      error
		flag     bool
		vIdx     uint
		itemList *riotmodel.ItemList
		items    map[string][]*riotmodel.ItemDTO // [ver-mode-lang]
	)

	vIdx, err = utils.ConvertVersionToIdx(version)
	if version == "" || err != nil {
		u.logger.Error("wrong version")
	}
	items = make(map[string][]*riotmodel.ItemDTO)
	key := "/items"

	ctx := context.Background()
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLangToLangStr(langCode)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/item.json",
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
		key2 := fmt.Sprintf("/items/%s", lang)
		for k, it := range itemList.Data {
			it.ID = k
			it.Description = utils.RemoveHTMLTags(it.Description)
			if it.Maps["11"] {
				tk := fmt.Sprintf("%d-classic-%s", vIdx, lang)
				items[tk] = append(items[tk], it)
			}
			if it.Maps["12"] {
				tk := fmt.Sprintf("%d-aram-%s", vIdx, lang)
				items[tk] = append(items[tk], it)
			}
			if it.Maps["30"] {
				tk := fmt.Sprintf("%d-cherry-%s", vIdx, lang)
				items[tk] = append(items[tk], it)
			}

			cmds = append(cmds, pipe.HSet(ctx, key2, fmt.Sprintf("%s@%d", k, vIdx), it))
		}
		if _, err = pipe.Exec(ctx); err != nil {
			u.logger.Error("redis store items failed", zap.Error(err))
		}
	}
	for k, its := range items {
		if buff, err = json.Marshal(its); err != nil {
			flag = true
			u.logger.Error("marshal item list failed", zap.Error(err))
			continue
		}
		if err = u.rdb.HSet(ctx, key, k, buff).Err(); err != nil {
			flag = true
			u.logger.Error("cache item list failed", zap.Error(err))
			continue
		}
	}
	if !flag {
		u.logger.Info("all " + version + "'s items fetch done")
	}
}

func (u *Updater) UpdatePerks(version string) {
	var (
		buff  []byte
		url   string
		err   error
		key   string
		vIdx  uint
		perks []*riotmodel.Perk
	)

	vIdx, err = utils.ConvertVersionToIdx(version)
	if version == "" || err != nil {
		u.logger.Error("wrong version")
	}
	ctx := context.Background()
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLangToLangStr(langCode)
		key = fmt.Sprintf("/perks/%s", lang)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		// fetch perk relation
		url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/runesReforged.json",
			version, lang)
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
			if err = u.rdb.HSet(ctx, key, fmt.Sprintf("%d@%d", p.ID, vIdx), p).Err(); err != nil {
				u.logger.Error("save perks failed", zap.Error(err))
				break
			}
		}

		// fetch perk detail (third-party data)
		// if lang == "en_US" {
		// 	lang = "default"
		// }
		// url = fmt.Sprintf("https://raw.communitydragon.org/%s/plugins/rcp-be-lol-game-data/global/%s/v1/perks.json",
		// 	version[:len(version)-2], strings.ToLower(lang))
		// if buff, err = u.fetcher.Get(url); err != nil {
		// 	u.logger.Error("get  perks failed", zap.Error(err))
		// 	continue
		// }
		// if err = json.Unmarshal(buff, &perkDetails); err != nil {
		// 	u.logger.Error("unmarshal json to item list failed", zap.Error(err))
		// 	continue
		// }
		// for _, p := range perkDetails {
		// 	if err = u.rdb.HSet(ctx, key, fmt.Sprintf("%d@%d", p.ID, vIdx), p).Err(); err != nil {
		// 		u.logger.Error("save perks failed", zap.Error(err))
		// 		break
		// 	}
		// }
	}
	u.logger.Info("all perk update done")
}

func (u *Updater) UpdateSpell(version string) {
	var (
		buff      []byte
		url       string
		err       error
		key       string
		vIdx      uint
		spellList *riotmodel.SpellList
	)

	vIdx, err = utils.ConvertVersionToIdx(version)
	if version == "" || err != nil {
		u.logger.Error("wrong version")
	}
	ctx := context.Background()
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLangToLangStr(langCode)
		key = "/spells/"
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		url = fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/summoner.json",
			version, lang)
		if buff, err = u.fetcher.Get(url); err != nil {
			u.logger.Error("get  spell failed", zap.Error(err))
			continue
		}
		if err = json.Unmarshal(buff, &spellList); err != nil {
			u.logger.Error("unmarshal json to item list failed", zap.Error(err))
			continue
		}

		if err = u.rdb.HSet(ctx, key, fmt.Sprintf("%d@%s", vIdx, lang), spellList).Err(); err != nil {
			u.logger.Error("redis store spell failed", zap.Error(err))
		}
	}
	u.logger.Info(fmt.Sprintf("all %s's summoner spell fetch done", version))
}

func isEnd(curVersion, endMark string) bool {
	cIdx, _ := utils.ConvertVersionToIdx(curVersion)
	eIdx, _ := utils.ConvertVersionToIdx(endMark)
	return cIdx >= eIdx
}
