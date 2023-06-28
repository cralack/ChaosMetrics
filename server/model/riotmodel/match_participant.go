package riotmodel

type ParticipantDto struct {
	MatchID string // 关联InfoDTO的外键

	Assists                        int      `json:"assists"`                        // 助攻数
	BaronKills                     int      `json:"baronKills"`                     // 击杀大龙数
	BountyLevel                    int      `json:"bountyLevel"`                    // 赏金等级
	ChampExperience                int      `json:"champExperience"`                // 英雄经验
	ChampLevel                     int      `json:"champLevel"`                     // 英雄等级
	ChampionId                     int      `json:"championId"`                     // 英雄ID
	ChampionName                   string   `json:"championName"`                   // 英雄名称
	ChampionTransform              int      `json:"championTransform"`              // 英雄转化（Kayn的转型）
	ConsumablesPurchased           int      `json:"consumablesPurchased"`           // 购买消耗品数
	DamageDealtToBuildings         int      `json:"damageDealtToBuildings"`         // 对建筑造成的伤害
	DamageDealtToObjectives        int      `json:"damageDealtToObjectives"`        // 对目标造成的伤害
	DamageDealtToTurrets           int      `json:"damageDealtToTurrets"`           // 对防御塔造成的伤害
	DamageSelfMitigated            int      `json:"damageSelfMitigated"`            // 减少的伤害
	Deaths                         int      `json:"deaths"`                         // 死亡数
	DetectorWardsPlaced            int      `json:"detectorWardsPlaced"`            // 放置控制视野守卫数
	DoubleKills                    int      `json:"doubleKills"`                    // 双杀数
	DragonKills                    int      `json:"dragonKills"`                    // 击杀小龙数
	FirstBloodAssist               bool     `json:"firstBloodAssist"`               // 是否协助一血
	FirstBloodKill                 bool     `json:"firstBloodKill"`                 // 是否击杀一血
	FirstTowerAssist               bool     `json:"firstTowerAssist"`               // 是否协助一塔
	FirstTowerKill                 bool     `json:"firstTowerKill"`                 // 是否击杀一塔
	GameEndedInEarlySurrender      bool     `json:"gameEndedInEarlySurrender"`      // 是否提前投降结束游戏
	GameEndedInSurrender           bool     `json:"gameEndedInSurrender"`           // 是否投降结束游戏
	GoldEarned                     int      `json:"goldEarned"`                     // 获得的金币
	GoldSpent                      int      `json:"goldSpent"`                      // 花费的金币
	IndividualPosition             string   `json:"individualPosition"`             // 个人位置
	InhibitorKills                 int      `json:"inhibitorKills"`                 // 击杀抑制塔数
	InhibitorTakedowns             int      `json:"inhibitorTakedowns"`             // 击杀抑制塔参与数
	InhibitorsLost                 int      `json:"inhibitorsLost"`                 // 失去的抑制塔数
	Item0                          int      `json:"item0"`                          // 物品0
	Item1                          int      `json:"item1"`                          // 物品1
	Item2                          int      `json:"item2"`                          // 物品2
	Item3                          int      `json:"item3"`                          // 物品3
	Item4                          int      `json:"item4"`                          // 物品4
	Item5                          int      `json:"item5"`                          // 物品5
	Item6                          int      `json:"item6"`                          // 物品6
	ItemsPurchased                 int      `json:"itemsPurchased"`                 // 购买物品数
	KillingSprees                  int      `json:"killingSprees"`                  // 连杀数
	Kills                          int      `json:"kills"`                          // 击杀数
	Lane                           string   `json:"lane"`                           // 路线
	LargestCriticalStrike          int      `json:"largestCriticalStrike"`          // 最大暴击伤害
	LargestKillingSpree            int      `json:"largestKillingSpree"`            // 最大连杀数
	LargestMultiKill               int      `json:"largestMultiKill"`               // 最大多杀数
	LongestTimeSpentLiving         int      `json:"longestTimeSpentLiving"`         // 存活最长时间
	MagicDamageDealt               int      `json:"magicDamageDealt"`               // 魔法伤害输出
	MagicDamageDealtToChampions    int      `json:"magicDamageDealtToChampions"`    // 对英雄造成的魔法伤害
	MagicDamageTaken               int      `json:"magicDamageTaken"`               // 承受的魔法伤害
	NeutralMinionsKilled           int      `json:"neutralMinionsKilled"`           // 击杀中立生物数
	NexusKills                     int      `json:"nexusKills"`                     // 击杀水晶数
	NexusTakedowns                 int      `json:"nexusTakedowns"`                 // 击杀水晶参与数
	NexusLost                      int      `json:"nexusLost"`                      // 失去的水晶数
	ObjectivesStolen               int      `json:"objectivesStolen"`               // 偷取目标数
	ObjectivesStolenAssists        int      `json:"objectivesStolenAssists"`        // 偷取目标参与数
	ParticipantId                  int      `json:"participantId"`                  // 参与者ID
	PentaKills                     int      `json:"pentaKills"`                     // 五杀数
	Perks                          PerksDto `json:"perks"`                          // 符文
	PhysicalDamageDealt            int      `json:"physicalDamageDealt"`            // 物理伤害输出
	PhysicalDamageDealtToChampions int      `json:"physicalDamageDealtToChampions"` // 对英雄造成的物理伤害
	PhysicalDamageTaken            int      `json:"physicalDamageTaken"`            // 承受的物理伤害
	ProfileIcon                    int      `json:"profileIcon"`                    // 头像图标ID
	Puuid                          string   `json:"puuid"`                          // 参与者UUID
	QuadraKills                    int      `json:"quadraKills"`                    // 四杀数
	RiotIdName                     string   `json:"riotIdName"`                     // Riot ID名称
	RiotIdTagline                  string   `json:"riotIdTagline"`                  // Riot ID标签
	Role                           string   `json:"role"`                           // 角色
	SightWardsBoughtInGame         int      `json:"sightWardsBoughtInGame"`         // 购买视野守卫数
	Spell1Casts                    int      `json:"spell1Casts"`                    // 技能1释放次数
	Spell2Casts                    int      `json:"spell2Casts"`                    // 技能2释放次数
	Spell3Casts                    int      `json:"spell3Casts"`                    // 技能3释放次数
	Spell4Casts                    int      `json:"spell4Casts"`                    // 技能4释放次数
	Summoner1Casts                 int      `json:"summoner1Casts"`                 // 召唤师技能1释放次数
	Summoner1Id                    int      `json:"summoner1Id"`                    // 召唤师技能1ID
	Summoner2Casts                 int      `json:"summoner2Casts"`                 // 召唤师技能2释放次数
	Summoner2Id                    int      `json:"summoner2Id"`                    // 召唤师技能2ID
	SummonerId                     string   `json:"summonerId"`                     // 召唤师ID
	SummonerLevel                  int      `json:"summonerLevel"`                  // 召唤师等级
	SummonerName                   string   `json:"summonerName"`                   // 召唤师名称
	TeamEarlySurrendered           bool     `json:"teamEarlySurrendered"`           // 队伍提前投降
	TeamId                         int      `json:"teamId"`                         // 队伍ID
	TeamPosition                   string   `json:"teamPosition"`                   // 队伍位置
	TimeCCingOthers                int      `json:"timeCCingOthers"`                // 控制敌方英雄时间
	TimePlayed                     int      `json:"timePlayed"`                     // 游戏时间
	TotalDamageDealt               int      `json:"totalDamageDealt"`               // 总伤害输出
	TotalDamageDealtToChampions    int      `json:"totalDamageDealtToChampions"`    // 对英雄造成的总伤害
	TotalDamageShieldedOnTeammates int      `json:"totalDamageShieldedOnTeammates"` // 对队友护盾总量
	TotalDamageTaken               int      `json:"totalDamageTaken"`               // 承受的总伤害
	TotalHeal                      int      `json:"totalHeal"`                      // 总治疗量
	TotalHealsOnTeammates          int      `json:"totalHealsOnTeammates"`          // 对队友的治疗量
	TotalMinionsKilled             int      `json:"totalMinionsKilled"`             // 总击杀小兵数
	TotalTimeCCDealt               int      `json:"totalTimeCCDealt"`               // 总控制时间
	TotalTimeSpentDead             int      `json:"totalTimeSpentDead"`             // 总死亡时间
	TotalUnitsHealed               int      `json:"totalUnitsHealed"`               // 总治疗单位数
	TripleKills                    int      `json:"tripleKills"`                    // 三杀数
	TrueDamageDealt                int      `json:"trueDamageDealt"`                // 真实伤害输出
	TrueDamageDealtToChampions     int      `json:"trueDamageDealtToChampions"`     // 对英雄造成的真实伤害
	TrueDamageTaken                int      `json:"trueDamageTaken"`                // 承受的真实伤害
	TurretKills                    int      `json:"turretKills"`                    // 击杀防御塔数
	TurretTakedowns                int      `json:"turretTakedowns"`                // 击杀防御塔参与数
	TurretsLost                    int      `json:"turretsLost"`                    // 失去的防御塔数
	UnrealKills                    int      `json:"unrealKills"`                    // 不可思议的击杀数
	VisionScore                    int      `json:"visionScore"`                    // 视野得分
	VisionWardsBoughtInGame        int      `json:"visionWardsBoughtInGame"`        // 购买视野守卫数
	WardsKilled                    int      `json:"wardsKilled"`                    // 摧毁守卫数
	WardsPlaced                    int      `json:"wardsPlaced"`                    // 放置守卫数
	Win                            bool     `json:"win"`                            // 是否获胜
}

type PerksDto struct {
	StatPerks *PerkStatsDto   `json:"statPerks"` // 符文属性
	Styles    []*PerkStyleDto `json:"styles"`    // 符文样式列表
}

type PerkStyleDto struct {
	Description string                   `json:"description"` // 描述
	Selections  []*PerkStyleSelectionDto `json:"selections"`  // 符文样式选择
	Style       int                      `json:"style"`       // 符文样式
}

type PerkStatsDto struct {
	Defense int `json:"defense"` // 防御属性
	Flex    int `json:"flex"`    // 弹性属性
	Offense int `json:"offense"` // 进攻属性
}

type PerkStyleSelectionDto struct {
	Perk int `json:"perk"` // 符文
	Var1 int `json:"var1"` // 变量1
	Var2 int `json:"var2"` // 变量2
	Var3 int `json:"var3"` // 变量3
}
