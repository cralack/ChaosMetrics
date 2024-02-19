package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const anaKey = "anamatch"

type Analyzer struct {
	logger    *zap.Logger
	db        *gorm.DB
	rdb       *redis.Client
	lock      *sync.Mutex
	scheduler scheduler.Scheduler
	Exit      chan struct{}

	curVersion    string
	analyzedCount []int64
	options       *options
	idxMap        map[int]string                    // idx->championName
	totalPlayed   map[string]int64                  // totalPlayed[version+loc+mode]
	banedCount    map[string]int                    // banedCount[name+ver+loc+mode]
	pickCount     map[string]int                    // pickCount[name+ver+loc+mode]
	chamTemplate  map[string]*riotmodel.ChampionDTO // chamTemplate[championName]
	itemMap       map[string]*riotmodel.ItemDTO     // itemMap[itemID@version]
	analyzed      map[string]*anres.ChampionDetail  // analyzed[name+ver+loc+mode]
	shoesList     map[uint]map[int]struct{}         // shoesList[version][shoes]
}

func NewAnalyzer(opts ...Option) *Analyzer {
	stgy := defaultOptions
	for _, opt := range opts {
		opt(stgy)
	}
	buff := global.ChaRDB.HGet(context.Background(), "/championlist", "idxmap").Val()
	idxMap := make(map[int]string)
	if err := json.Unmarshal([]byte(buff), &idxMap); err != nil {
		global.ChaLogger.Error("failed to parse idx map", zap.Error(err))
	}

	return &Analyzer{
		logger:        global.ChaLogger,
		db:            global.ChaDB,
		rdb:           global.ChaRDB,
		lock:          &sync.Mutex{},
		scheduler:     scheduler.NewSchdule(),
		Exit:          make(chan struct{}),
		curVersion:    global.ChaRDB.HGet(context.Background(), "/version", "cur").Val(),
		analyzedCount: make([]int64, 16),
		options:       stgy,
		idxMap:        idxMap,
		totalPlayed:   make(map[string]int64),
		banedCount:    make(map[string]int),
		pickCount:     make(map[string]int),
		chamTemplate:  make(map[string]*riotmodel.ChampionDTO),
		itemMap:       make(map[string]*riotmodel.ItemDTO),
		analyzed:      make(map[string]*anres.ChampionDetail),
		shoesList:     make(map[uint]map[int]struct{}),
	}
}

func (a *Analyzer) Analyze() {
	go a.scheduler.Schedule()
	go a.handleMatches()

	a.loadChampionTemplate()
	for _, loc := range a.options.Loc {
		a.loadMatch(loc)
	}
	<-a.Exit
	defer a.store()
}

func (a *Analyzer) loadItem(itemId int, version uint) (res *riotmodel.ItemDTO) {
	var (
		err    error
		has    bool
		buffer string
	)
	key := fmt.Sprintf("%d@%d", itemId, version)
	if res, has = a.itemMap[key]; has {
		return res
	}
	if buffer, err = a.rdb.HGet(context.Background(), "/items/zh_CN", key).Result(); err != nil {
		a.logger.Error("load item from redis failed", zap.Error(err))
		return nil
	}
	if err = json.Unmarshal([]byte(buffer), &res); err != nil {
		a.logger.Error("wrong item buffer from redis")
		return nil
	}
	if itemId == 1001 {
		if _, has = a.shoesList[version]; !has {
			a.shoesList[version] = make(map[int]struct{})
		}
		shoesL := res.Into
		for _, shoe := range shoesL {
			shoeId, err := strconv.Atoi(shoe)
			if err != nil {
				a.logger.Error("get shoe failed", zap.Error(err))
			}
			a.shoesList[version][shoeId] = struct{}{}
		}
	}
	a.itemMap[key] = res
	return
}

func (a *Analyzer) loadMatch(loCode riotmodel.LOCATION) {
	var (
		err     error
		matches []*riotmodel.MatchDB
	)

	loc, _ := utils.ConvertLocationToLoHoSTR(loCode)

	// count analyzed
	var totalCount int64
	if err = a.db.Model(&riotmodel.MatchDB{}).Where("analyzed = ?",
		true).Count(&totalCount).Error; err != nil {
		a.logger.Error("count analyzed failed", zap.Error(err))
		return
	}
	a.lock.Lock()
	a.analyzedCount[loCode] += totalCount
	a.lock.Unlock()

	// count unanalyzed
	if err = a.db.Model(&riotmodel.MatchDB{}).Where("analyzed = ?",
		false).Count(&totalCount).Error; err != nil {
		a.logger.Error("count analyzed failed", zap.Error(err))
		return
	}
	// start counter
	go a.counter(int(totalCount), loCode)

	// chunk
	totalSize := int(totalCount)
	chunkSize := a.options.BatchSize
	for i := 0; i < totalSize; i += chunkSize {
		end := i + chunkSize
		if end > totalSize {
			end = totalSize
		}
		matches = make([]*riotmodel.MatchDB, 0, a.options.BatchSize)
		if err = a.db.Offset(i).Limit(chunkSize).Where("loc = ?", loc).Where("analyzed = ?",
			false).Where("analyzed = ?", false).Preload(
			"Participants").Find(&matches).Error; err != nil {
			a.logger.Error("load match failed", zap.Error(err))
		}
		a.scheduler.Push(&scheduler.Task{
			Type: anaKey,
			Loc:  loc,
			Data: matches,
		})
	}

	a.scheduler.Push(&scheduler.Task{
		Type: "finish",
		Data: nil,
	})
	return
}

// prove champion info
func (a *Analyzer) loadChampionTemplate() {
	var (
		err           error
		ctx           context.Context
		vIdx          uint
		values        []interface{}
		championNames []string
		keys          []string
	)

	if vIdx, err = utils.ConvertVersionToIdx(a.curVersion); err != nil {
		a.logger.Error("wrong version", zap.Error(err))
	}
	ctx = context.Background()
	curLang := utils.ConvertLangToLangStr(riotmodel.LANG_zh_CN)
	key := "/champions/" + curLang
	// get champion name list
	buff := a.rdb.HGet(context.Background(), "/championlist", strconv.Itoa(int(vIdx))).Val()
	if err = json.Unmarshal([]byte(buff), &championNames); err != nil {
		a.logger.Error("unmarshal champion list failed", zap.Error(err))
	}
	keys = make([]string, len(championNames))
	for i, k := range championNames {
		keys[i] = fmt.Sprintf("%s@%d", k, vIdx)
	}
	// get champions by name list
	values, err = a.rdb.HMGet(ctx, key, keys...).Result()
	if err != nil {
		a.logger.Error("get champions failed", zap.Error(err))
	}

	// parse
	for _, v := range values {
		var cham *riotmodel.ChampionDTO
		if err = json.Unmarshal([]byte(v.(string)), &cham); err != nil {
			a.logger.Error("unmarshal champion failed", zap.Error(err))
		} else {
			// assign to map
			a.chamTemplate[strings.ToLower(cham.ID)] = cham
		}
	}
	a.logger.Debug("load champion succeed", zap.Int("len", len(a.chamTemplate)))
	return
}

func (a *Analyzer) handleMatches() {
	for {
		req := a.scheduler.Pull()
		switch req.Type {
		case "finish":
			a.Exit <- struct{}{}
			return
		case anaKey:
			matches := req.Data.([]*riotmodel.MatchDB)
			for _, match := range matches {
				a.AnalyzeSingleMatch(match)
			}
		}
	}

}

func (a *Analyzer) AnalyzeSingleMatch(match *riotmodel.MatchDB) {
	if match.GameMode == "CHERRY" || len(match.Participants) == 0 {
		return
	}
	var (
		has bool
		err error
	)
	// get param
	loCode := utils.ConvertLocStrToLocation(match.Loc)
	curVersion := match.GameVersion

	if err != nil {
		a.logger.Error("wrong match version")
		return
	}

	// count ban rate
	if match.GameMode == "CLASSIC" {
		var bans []int
		if err = json.Unmarshal([]byte(match.Bans), &bans); err != nil {
			a.logger.Error("wrong bans", zap.Error(err))
		}
		for _, id := range bans {
			if id == -1 {
				continue
			}
			k := GetID(a.idxMap[id], curVersion, match.Loc, match.GameMode)
			a.banedCount[k]++
		}
	}

	// match partic
	for _, par := range match.Participants {
		// skip BOT game
		if par.Puuid == "BOT" {
			return
		}
		keyName := par.ChampionName
		var (
			// chamIdx  int
			tarId    string
			tmp      *anres.ChampionDetail
			template *riotmodel.ChampionDTO
		)
		// get champion data template
		if template, has = a.chamTemplate[strings.ToLower(keyName)]; !has {
			a.logger.Error(keyName + "doesnt exist")
			return
		}

		// match champion && version && loc && gamemode
		tarId = GetID(keyName, curVersion, match.Loc, match.GameMode)
		if tmp, has = a.analyzed[tarId]; !has {
			tmp = &anres.ChampionDetail{
				Loc:      match.Loc,
				Version:  curVersion,
				MetaName: keyName,
				Key:      template.Key,
				Name:     template.Name,
				Title:    template.Title,
				Image:    template.Image,
				GameMode: match.GameMode,
				ItemWin: map[string]map[string]int{
					"fir": {},
					"tri": {},
					"oth": {},
					"sho": {},
				},
				PerkWin:  make(map[string]int),
				SkillWin: make(map[string]int),
				SpellWin: make(map[string]int),
			}
			tmp.ID = tarId
		}
		a.pickCount[tarId]++
		// analyze data
		if err = a.handleAnares(tmp, par); err != nil {
			a.logger.Error("parse match failed", zap.Error(err))
			return
		}
		a.analyzed[tarId] = tmp

	}

	a.lock.Lock()
	defer a.lock.Unlock()
	// analyzed match count
	a.totalPlayed[GetID("", curVersion, match.Loc, match.GameMode)]++
	a.analyzedCount[loCode]++
	return
}

func (a *Analyzer) store() {
	analyzedDetail := make([]*anres.ChampionDetail, 0, len(a.analyzed))
	analyzed := make(map[string][]*anres.ChampionBrief)
	cmd := make([]*redis.IntCmd, 0, len(a.analyzed))
	ctx := context.Background()
	pipe := a.rdb.Pipeline()

	// traverse all analyzed result
	for key, cham := range a.analyzed {
		k := GetID("", cham.Version, cham.Loc, cham.GameMode)
		cham.BanRate = float32(a.banedCount[key]) / float32(a.totalPlayed[k])
		cham.PickRate = float32(a.pickCount[key]) / float32(a.totalPlayed[k])

		analyzedDetail = append(analyzedDetail, cham)
		cmd = append(cmd, pipe.HSet(ctx, "/champion_detail", cham.ID, cham))
		// store analyzed brief
		vidx, _ := utils.ConvertVersionToIdx(cham.Version)
		idx := fmt.Sprintf("%d_%s@%s", vidx, cham.GameMode, cham.Loc)
		if _, has := analyzed[idx]; !has {
			analyzed[idx] = make([]*anres.ChampionBrief, 0, 200)
		}
		analyzed[idx] = append(analyzed[idx], &anres.ChampionBrief{
			Image:          cham.Image,
			MetaName:       cham.MetaName,
			WinRate:        cham.WinRate,
			PickRate:       cham.PickRate,
			BanRate:        cham.BanRate,
			AvgDamageDealt: cham.AvgDamageDealt,
			AvgDeadTime:    cham.AvgDeadTime,
		})
	}
	// store detail
	if _, err := pipe.Exec(ctx); err != nil {
		a.logger.Error("store analyzed result to redis failed", zap.Error(err))
	}

	// store brief
	for idx, list := range analyzed {
		sort.Slice(list, func(i, j int) bool {
			return list[i].MetaName < list[j].MetaName
		})
		buff, err := json.Marshal(list)
		if err != nil {
			a.logger.Error("marshal champion brief failed", zap.Error(err))
		}
		a.rdb.HSet(ctx, "/champion_brief", idx, buff)
	}
	a.rdb.HSet(ctx, "/lastupdate", "analyzer", time.Now().Unix())
}

func (a *Analyzer) handleAnares(tar *anres.ChampionDetail, par *riotmodel.ParticipantDB) error {
	var (
		// buff         []byte
		verIdx       uint
		startCapital int
		startItem    []string
		triItems     []string
		otheItems    []string
		shoe         string
		skillBuild   []byte
		spell        string
		err          error
		judge        Judger
		item         *riotmodel.ItemDTO
	)
	// init val
	if verIdx, err = utils.ConvertVersionToIdx(tar.Version); err != nil {
		a.logger.Error("wrong version")
		return err
	}

	switch tar.GameMode {
	case "ARAM":
		startCapital = 1400
		judge = JudgeARAM()
	case "CLASSIC":
		startCapital = 500
		judge = JudgeClassic()
	}
	// parse item build
	if err = json.Unmarshal([]byte(par.Build.Item), &par.Build.ItemOrder); err != nil {
		return errors.New("unmarshal item failed" + err.Error())
	}
	spendMoney := 0
	if a.shoesList[verIdx] == nil {
		a.loadItem(1001, verIdx)
	}
	for _, it := range par.Build.ItemOrder {
		item = a.loadItem(it.ItemID, verIdx)
		// startItems
		if item.Gold.Total != 0 && spendMoney+item.Gold.Total <= startCapital {
			spendMoney += item.Gold.Total
			startItem = append(startItem, strconv.Itoa(it.ItemID))
			continue
		}
		// triItems
		if item.Depth == 3 && len(triItems) < 3 {
			triItems = append(triItems, strconv.Itoa(it.ItemID))
			continue
		}
		// otherItems
		if item.Depth == 3 && len(triItems) == 3 {
			otheItems = append(otheItems, strconv.Itoa(it.ItemID))
			continue
		}
		// shoe
		if _, has := a.shoesList[verIdx][it.ItemID]; has {
			shoe = strconv.Itoa(it.ItemID)
		}
	}
	// parse skill build
	if err = json.Unmarshal([]byte(par.Build.Skill), &par.Build.SkillOrder); err != nil {
		return errors.New("unmarshal skill failed" + err.Error())
	}
	skiMap := map[int]byte{1: 'Q', 2: 'W', 3: 'E'}
	ski := make([]int, 5)
	for _, s := range par.Build.SkillOrder {
		ski[s]++
		if ski[s] == 5 {
			skillBuild = append(skillBuild, skiMap[s])
		}
	}
	// summoner spell
	if par.Summoner1Id < par.Summoner2Id {
		spell = fmt.Sprintf("%d,%d", par.Summoner1Id, par.Summoner2Id)
	} else {
		spell = fmt.Sprintf("%d,%d", par.Summoner2Id, par.Summoner1Id)
	}

	// count
	if par.Win {
		tar.TotalWin++
		tar.PerkWin[par.Build.Perk]++
		tar.ItemWin["fir"][strings.Join(startItem, ",")]++
		tar.ItemWin["tri"][strings.Join(triItems, ",")]++
		for _, oth := range otheItems {
			tar.ItemWin["oth"][oth]++
		}
		tar.ItemWin["sho"][shoe]++
		if len(skillBuild) > 0 {
			tar.SkillWin[string(skillBuild)]++
		}
		tar.SpellWin[spell]++
	}

	tar.WinRate = tar.TotalWin / (tar.TotalPlayed + 1)
	tar.AvgKDA = (tar.AvgKDA*tar.TotalPlayed + par.KDA) / (tar.TotalPlayed + 1)
	tar.AvgKP = (tar.AvgKP*tar.TotalPlayed + par.KP) / (tar.TotalPlayed + 1)
	tar.AvgDamageDealt = (tar.AvgDamageDealt*tar.TotalPlayed + float32(par.DamageDealt)) / (tar.TotalPlayed + 1)
	tar.AvgDamageTaken = (tar.AvgDamageTaken*tar.TotalPlayed + float32(par.DamageToken)) / (tar.TotalPlayed + 1)
	tar.AvgTimeCCing = (tar.AvgTimeCCing*tar.TotalPlayed + float32(par.TimeCCingOthers)) / (tar.TotalPlayed + 1)
	tar.AvgVisionScore = (tar.AvgVisionScore*tar.TotalPlayed + float32(par.VisionScore)) / (tar.TotalPlayed + 1)
	tar.AvgDeadTime = (tar.AvgDeadTime*tar.TotalPlayed + float32(par.TotalTimeSpentDead)) / (tar.TotalPlayed + 1)
	tar.TotalPlayed++
	tar.RankScore = 0

	judge(tar)
	return nil
}

func (a *Analyzer) counter(total int, loc riotmodel.LOCATION) {
	var (
		cur  int64
		rate float32
	)
	ticker := time.NewTicker(time.Second * 5)

	for {
		// count
		cur = a.analyzedCount[loc]
		rate = float32(cur) / float32(total)
		a.logger.Info(fmt.Sprintf("analyzed %05.02f%% (%d/%d) match", rate*100, cur, total))
		// store
		// a.store()
		<-ticker.C
	}
}

func GetID(name, version, loc, mode string) string {
	vidx, _ := utils.ConvertVersionToIdx(version)
	if name == "" {
		return fmt.Sprintf("%d-%s-%s",
			vidx,
			strings.ToLower(loc),
			strings.ToLower(mode),
		)
	}

	return fmt.Sprintf("%s-%d-%s-%s",
		strings.ToLower(name),
		vidx,
		strings.ToLower(loc),
		strings.ToLower(mode),
	)
}
