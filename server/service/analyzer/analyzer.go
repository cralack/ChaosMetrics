package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/updater"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const anaKey = "anamatch"

type Analyzer struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.Client
	lock   *sync.Mutex
	schd   scheduler.Scheduler
	updt   *updater.Updater

	curVersion    string
	analyzedCount []int64
	stgy          *Strategy
	chamTemplate  map[string]*riotmodel.ChampionDTO // chamTemplate[championName]
	itemMap       map[string]*riotmodel.ItemDTO     // itemMap[itemID@version]
	analyzed      map[uint]*anres.Champion          // analyzed[chamId+verId+loc+mode]
	shoesList     map[uint]map[int]struct{}         // shoesList[version][shoes]
}

func NewAnalyzer(opts ...Option) *Analyzer {
	stgy := defaultStrategy
	for _, opt := range opts {
		opt(stgy)
	}

	return &Analyzer{
		logger:        global.GVA_LOG,
		db:            global.GVA_DB,
		rdb:           global.GVA_RDB,
		lock:          &sync.Mutex{},
		schd:          scheduler.NewSchdule(),
		updt:          updater.NewRiotUpdater(),
		analyzedCount: make([]int64, 16),
		stgy:          stgy,
		// matchMap:    make(map[string]map[string]*riotmodel.MatchDTO),
		chamTemplate: make(map[string]*riotmodel.ChampionDTO),
		itemMap:      make(map[string]*riotmodel.ItemDTO),
		analyzed:     make(map[uint]*anres.Champion),
		shoesList:    make(map[uint]map[int]struct{}),
	}
}

func (a *Analyzer) Analyze() {
	exit := make(chan struct{})
	go a.schd.Schedule()
	go a.handleMatches(exit)

	a.loadChampionTemplate()
	for _, loc := range a.stgy.Loc {
		a.loadMatch(loc)
	}
	<-exit
	defer a.store()
	// matchId := "TW2_81882122"
	// a.AnalyzeSingleMatch(matchId)
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
		if _, has := a.shoesList[version]; !has {
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

func (a *Analyzer) loadMatch(loCode uint) {
	var (
		err     error
		matches []*riotmodel.MatchDB
	)

	loc, _ := utils.ConvertHostURL(loCode)
	// preload when begin
	if err = a.db.Where("loc = ?", loc).Where("analyzed = ?", false).Preload(
		"Participants").Find(&matches).Error; err != nil {
		// or not
		// if err = a.db.Where("loc = ?", loc).Where("analyzed = ?",
		// 	false).Find(&matches).Error; err != nil {
		a.logger.Error("load match from gorm db failed", zap.Error(err))
	}

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

	// start counter
	go a.counter(len(matches), loCode)

	// chunk if oversize
	if len(matches) > a.stgy.BatchSize {
		totalSize := len(matches)
		chunkSize := a.stgy.BatchSize
		for i := 0; i < totalSize; i += chunkSize {
			end := i + chunkSize
			if end > totalSize {
				end = totalSize
			}
			tmp := matches[i:end]
			// preload and send matches
			// if err := a.db.Preload("Participants").Find(&tmp).Error; err != nil {
			// 	a.logger.Error("load participant failed", zap.Error(err))
			// }
			a.schd.Push(&scheduler.Task{
				Type: anaKey,
				Loc:  loc,
				Data: tmp,
			})
		}
	} else { // !oversize
		// if err := a.db.Preload("Participants").Find(&matches).Error; err != nil {
		// 	a.logger.Error("load participant failed", zap.Error(err))
		// }
		a.schd.Push(&scheduler.Task{
			Type: anaKey,
			Loc:  loc,
			Data: matches,
		})
	}

	a.schd.Push(&scheduler.Task{
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
		championNames []string
		keys          []string
	)
	curVersion := "13.16.1"
	if vIdx, err = utils.ConvertVersionToIdx(curVersion); err != nil {
		a.logger.Error("wrong version", zap.Error(err))
	}
	ctx = context.Background()
	curLang := utils.ConvertLanguageCode(riotmodel.LANG_zh_CN)
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
	values, err := a.rdb.HMGet(ctx, key, keys...).Result()
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
	a.logger.Debug("succeed", zap.Int("championmap", len(a.chamTemplate)))
	return
}

func (a *Analyzer) handleMatches(exit chan struct{}) {
	for {
		req := a.schd.Pull()
		switch req.Type {
		case "finish":
			exit <- struct{}{}
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
		// tar     []*anres.Champion
		has     bool
		verIdx  uint
		modeIdx uint
		err     error
	)
	// get param
	loCode := utils.ConverHostLoCode(match.Loc)
	curVersion := match.GameVersion
	verIdx, err = utils.ConvertVersionToIdx(curVersion)
	if err != nil {
		a.logger.Error("wrong match version")
		return
	}
	switch match.GameMode {
	case "ARAM":
		modeIdx = riotmodel.ARAM
	case "CHERRY":
		modeIdx = riotmodel.CHERRY
	case "CLASSIC":
		modeIdx = riotmodel.CLASSIC
	case "CONVERGENCE":
		modeIdx = riotmodel.CONVERGENCE
	case "NEXUSBLITZ":
		modeIdx = riotmodel.NEXUSBLITZ
	}
	// init champion each version?
	// if match.GameMode == "CLASSIC" {
	// 	var bans []int
	// 	if err = json.Unmarshal([]byte(match.Bans), &bans); err != nil {
	// 		a.logger.Error("wrong bans", zap.Error(err))
	// 	}
	// 	for _, id := range bans {
	// 		banId := uint(id)*1e8 + verIdx*1e4 + loCode*1e2 + modeIdx
	// 	}
	// }

	// match partic
	// tar = make([]*anres.Champion, 0, len(match.Participants))
	for _, par := range match.Participants {
		// skip BOT game
		if par.Puuid == "BOT" {
			return
		}
		keyName := par.ChampionName
		var (
			chamIdx  int
			tarId    uint
			tmp      *anres.Champion
			template *riotmodel.ChampionDTO
		)
		// get champion data template
		if cham, has := a.chamTemplate[strings.ToLower(keyName)]; !has {
			a.logger.Error(keyName + "doesnt exist")
			return
		} else {
			chamIdx, err = strconv.Atoi(cham.Key)
			if err != nil {
				a.logger.Error("wrong key")
				return
			}
			template = cham
		}
		if template == nil {
			return
		}
		// match champion && version && loc && gamemode
		tarId = uint(chamIdx)*1e8 + verIdx*1e4 + loCode*1e2 + modeIdx
		if tmp, has = a.analyzed[tarId]; !has {
			tmp = &anres.Champion{
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
	a.analyzedCount[loCode]++
	return
}

func (a *Analyzer) counter(total int, loc uint) {
	var (
		cur  int64
		rate float32
	)
	ticker := time.NewTicker(time.Second * 15)

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

func (a *Analyzer) store() {
	analyzed := make([]*anres.Champion, 0, len(a.analyzed))
	cmd := make([]*redis.IntCmd, 0, len(a.analyzed))
	ctx := context.Background()
	pipe := a.rdb.Pipeline()
	key := "/analyzed"
	for _, cham := range a.analyzed {
		analyzed = append(analyzed, cham)
		cmd = append(cmd, pipe.HSet(ctx, key, cham.ID, cham))
	}
	// store
	if _, err := pipe.Exec(ctx); err != nil {
		a.logger.Error("store analyzed result to redis failed")
	}
	if err := a.db.CreateInBatches(analyzed, 100).Error; err != nil {
		a.logger.Error("store analyzed result to db failed", zap.Error(err))
	}
}

func (a *Analyzer) handleAnares(tar *anres.Champion, par *riotmodel.ParticipantDB) error {
	var (
		buff         []byte
		totalCount   int64
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
	totalCount = a.analyzedCount[utils.ConverHostLoCode(tar.Loc)] + 1
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
	if buff, err = json.Marshal(tar.PerkWin); err != nil {
		return errors.New("marshal perks faild" + err.Error())
	} else {
		tar.PerkSTR = string(buff)
	}
	if buff, err = json.Marshal(tar.ItemWin); err != nil {
		return errors.New("marshal perks faild" + err.Error())
	} else {
		tar.ItemSTR = string(buff)
	}
	if buff, err = json.Marshal(tar.SkillWin); err != nil {
		return errors.New("marshal skill faild" + err.Error())
	} else {
		tar.SkillSTR = string(buff)
	}
	if buff, err = json.Marshal(tar.SpellWin); err != nil {
		return errors.New("marshal spell faild" + err.Error())
	} else {
		tar.SpellSTR = string(buff)
	}

	tar.WinRate = tar.TotalWin / (tar.TotalPlayed + 1)
	tar.PickRate = (tar.TotalPlayed + 1) / float32(totalCount)
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
