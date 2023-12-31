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
	if buff, err = u.fetcher.Get(fetcher.NewTask(
		fetcher.WithURL(url),
	)); err != nil || buff == nil {
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

	// get chamlist from rdb or riot
	if buffer = u.rdb.HGet(ctx, "/championlist", fmt.Sprintf("%dzh_CN", vIdx)).Val(); buffer == "" {
		// get champion chamList
		url = fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", version)
		if buff, err = u.fetcher.Get(fetcher.NewTask(
			fetcher.WithURL(url),
		)); err != nil || buff == nil {
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
	}

	// update champion for each lang
	for _, langCode := range u.stgy.Lang {
		lang := utils.ConvertLanguageCode(langCode)
		key := fmt.Sprintf("/champions/%s", lang)
		u.rdb.Expire(ctx, key, u.stgy.LifeTime)
		cmds := make([]*redis.IntCmd, 0, len(cList))
		pipe := u.rdb.Pipeline()
		cnt := 0

		for _, chamID := range cList {
			if buffer = u.rdb.HGet(ctx, key, fmt.Sprintf("%s@%d", chamID, vIdx)).Val(); buffer != "" {
				continue
			}

			// fetch buffer
			url = fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json",
				version, lang, chamID)
			if buff, err = u.fetcher.Get(fetcher.NewTask(
				fetcher.WithURL(url),
			)); err != nil || buff == nil {
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
				// u.logger.Debug(fmt.Sprintf("fetch %03d/%03d %s@%s succeed",
				// 	cnt, len(cList), chamID, lang))
			}
			cham = tmp.Data[chamID]
			cmds = append(cmds, pipe.HSet(ctx, key, fmt.Sprintf("%s@%d", cham.ID, vIdx), cham))
		}

		// store chamlist
		if _, err = pipe.Exec(ctx); err != nil {
			u.logger.Error("redis store champions failed", zap.Error(err))
		} else {
			// champion data's save status
			if buff, err = json.Marshal(cList); err != nil {
				u.logger.Error("marshal champion list failed", zap.Error(err))
			}
			if err = u.rdb.HSet(ctx, "/championlist", fmt.Sprintf("%d%s", vIdx, lang), buff).Err(); err != nil {
				u.logger.Error("failed", zap.Error(err))
			} else {
				u.logger.Info(fmt.Sprintf("%d@%s's champion(%d/%d) store done", vIdx, lang, cnt, len(cList)))
			}
		}
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
		if buff, err = u.fetcher.Get(fetcher.NewTask(
			fetcher.WithURL(url),
		)); err != nil {
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
		u.logger.Info("all " + version + "'s items fetch done")
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
		if buff, err = u.fetcher.Get(fetcher.NewTask(
			fetcher.WithURL(url),
		)); err != nil {
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
		if buff, err = u.fetcher.Get(fetcher.NewTask(
			fetcher.WithURL(url),
		)); err != nil {
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
