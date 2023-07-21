package riotmodel

import (
	"gorm.io/gorm"
)

type ParticipantDto struct {
	gorm.Model
	MatchID     uint   `gorm:"column:match_id"`                              // 匹配matchDTO_id
	MetaMatchID string `gorm:"column:meta_match_id;index;type:varchar(100)"` // 比赛ID
	
	Assists                        int       `json:"assists" gorm:"column:assists;type:varchar(100)"`                                    // 助攻数
	BaronKills                     int       `json:"baronKills" gorm:"column:baron_kills"`                                               // 击杀大龙数
	BountyLevel                    int       `json:"bountyLevel" gorm:"column:bounty_level;type:smallint"`                               // 赏金等级
	ChampExperience                int       `json:"champExperience" gorm:"column:champ_experience;type:int"`                            // 英雄经验
	ChampLevel                     int       `json:"champLevel" gorm:"column:champ_level;type:smallint"`                                 // 英雄等级
	ChampionId                     int       `json:"championId" gorm:"column:champion_id;type:smallint"`                                 // 英雄ID
	ChampionName                   string    `json:"championName" gorm:"column:champion_name;type:varchar(100)"`                         // 英雄名称
	ChampionTransform              int       `json:"championTransform" gorm:"column:champion_transform;type:smallint"`                   // 英雄转化（Kayn的转型）
	ConsumablesPurchased           int       `json:"consumablesPurchased" gorm:"column:consumables_purchased;type:smallint"`             // 购买消耗品数
	DamageDealtToBuildings         int       `json:"damageDealtToBuildings" gorm:"column:damage_dealt_to_buildings;type:int"`            // 对建筑造成的伤害
	DamageDealtToObjectives        int       `json:"damageDealtToObjectives" gorm:"column:damage_dealt_to_objectives;type:int"`          // 对目标造成的伤害
	DamageDealtToTurrets           int       `json:"damageDealtToTurrets" gorm:"column:damage_dealt_to_turrets;type:int"`                // 对防御塔造成的伤害
	DamageSelfMitigated            int       `json:"damageSelfMitigated" gorm:"column:damage_self_mitigated;type:int"`                   // 减少的伤害
	Deaths                         int       `json:"deaths" gorm:"column:deaths;type:smallint"`                                          // 死亡数
	DetectorWardsPlaced            int       `json:"detectorWardsPlaced" gorm:"column:detector_wards_placed;type:smallint"`              // 放置控制视野守卫数
	DoubleKills                    int       `json:"doubleKills" gorm:"column:double_kills;type:smallint"`                               // 双杀数
	DragonKills                    int       `json:"dragonKills" gorm:"column:dragon_kills;type:smallint"`                               // 击杀小龙数
	FirstBloodAssist               bool      `json:"firstBloodAssist" gorm:"column:first_blood_assist"`                                  // 是否协助一血
	FirstBloodKill                 bool      `json:"firstBloodKill" gorm:"column:first_blood_kill"`                                      // 是否击杀一血
	FirstTowerAssist               bool      `json:"firstTowerAssist" gorm:"column:first_tower_assist"`                                  // 是否协助一塔
	FirstTowerKill                 bool      `json:"firstTowerKill" gorm:"column:first_tower_kill"`                                      // 是否击杀一塔
	GameEndedInEarlySurrender      bool      `json:"gameEndedInEarlySurrender" gorm:"column:game_ended_in_early_surrender"`              // 是否提前投降结束游戏
	GameEndedInSurrender           bool      `json:"gameEndedInSurrender" gorm:"column:game_ended_in_surrender"`                         // 是否投降结束游戏
	GoldEarned                     int       `json:"goldEarned" gorm:"column:gold_earned;type:int"`                                      // 获得的金币
	GoldSpent                      int       `json:"goldSpent" gorm:"column:gold_spent;type:int"`                                        // 花费的金币
	IndividualPosition             string    `json:"individualPosition" gorm:"column:individual_position;type:varchar(100)"`             // 个人位置
	InhibitorKills                 int       `json:"inhibitorKills" gorm:"column:inhibitor_kills;type:smallint"`                         // 击杀抑制塔数
	InhibitorTakedowns             int       `json:"inhibitorTakedowns" gorm:"column:inhibitor_takedowns;type:smallint"`                 // 击杀抑制塔参与数
	InhibitorsLost                 int       `json:"inhibitorsLost" gorm:"column:inhibitors_lost;type:smallint"`                         // 失去的抑制塔数
	Item0                          int       `json:"item0" gorm:"column:item0;type:int"`                                                 // 物品0
	Item1                          int       `json:"item1" gorm:"column:item1;type:int"`                                                 // 物品1
	Item2                          int       `json:"item2" gorm:"column:item2;type:int"`                                                 // 物品2
	Item3                          int       `json:"item3" gorm:"column:item3;type:int"`                                                 // 物品3
	Item4                          int       `json:"item4" gorm:"column:item4;type:int"`                                                 // 物品4
	Item5                          int       `json:"item5" gorm:"column:item5;type:int"`                                                 // 物品5
	Item6                          int       `json:"item6" gorm:"column:item6;type:int"`                                                 // 物品6
	ItemsPurchased                 int       `json:"itemsPurchased" gorm:"column:items_purchased;type:smallint"`                         // 购买物品数
	KillingSprees                  int       `json:"killingSprees" gorm:"column:killing_sprees;type:smallint"`                           // 连杀数
	Kills                          int       `json:"kills" gorm:"column:kills;type:smallint"`                                            // 击杀数
	Lane                           string    `json:"lane" gorm:"column:lane;type:varchar(100)"`                                          // 路线
	LargestCriticalStrike          int       `json:"largestCriticalStrike" gorm:"column:largest_critical_strike;type:smallint"`          // 最大暴击伤害
	LargestKillingSpree            int       `json:"largestKillingSpree" gorm:"column:largest_killing_spree;type:smallint"`              // 最大连杀数
	LargestMultiKill               int       `json:"largestMultiKill" gorm:"column:largest_multi_kill;type:smallint"`                    // 最大多杀数
	LongestTimeSpentLiving         int       `json:"longestTimeSpentLiving" gorm:"column:longest_time_spent_living;type:smallint"`       // 存活最长时间
	MagicDamageDealt               int       `json:"magicDamageDealt" gorm:"column:magic_damage_dealt;type:int"`                         // 魔法伤害输出
	MagicDamageDealtToChampions    int       `json:"magicDamageDealtToChampions" gorm:"column:magic_damage_dealt_to_champions;type:int"` // 对英雄造成的魔法伤害
	MagicDamageTaken               int       `json:"magicDamageTaken" gorm:"column:magic_damage_taken;type:int"`                         // 承受的魔法伤害
	NeutralMinionsKilled           int       `json:"neutralMinionsKilled" gorm:"column:neutral_minions_killed;type:smallint"`            // 击杀中立生物数
	NexusKills                     int       `json:"nexusKills" gorm:"column:nexus_kills;type:smallint"`                                 // 击杀水晶数
	NexusTakedowns                 int       `json:"nexusTakedowns" gorm:"column:nexus_takedowns;type:smallint"`                         // 击杀水晶参与数
	NexusLost                      int       `json:"nexusLost" gorm:"column:nexus_lost;type:smallint"`                                   // 失去的水晶数
	ObjectivesStolen               int       `json:"objectivesStolen" gorm:"column:objectives_stolen;type:smallint"`                     // 偷取目标数
	ObjectivesStolenAssists        int       `json:"objectivesStolenAssists" gorm:"column:objectives_stolen_assists;type:smallint"`      // 偷取目标参与数
	ParticipantId                  int       `json:"participantId" gorm:"column:participant_id;type:smallint"`                           // 参与者ID
	PentaKills                     int       `json:"pentaKills" gorm:"column:penta_kills;type:smallint"`                                 // 五杀数
	PerksMeta                      *PerksDto `json:"perks" gorm:"-"`                                                                     // 符文
	PerksORM                       *PerksORM `gorm:"embedded;embeddedPrefix:perks_"`
	PhysicalDamageDealt            int       `json:"physicalDamageDealt" gorm:"column:physical_damage_dealt;type:int"`                         // 物理伤害输出
	PhysicalDamageDealtToChampions int       `json:"physicalDamageDealtToChampions" gorm:"column:physical_damage_dealt_to_champions;type:int"` // 对英雄造成的物理伤害
	PhysicalDamageTaken            int       `json:"physicalDamageTaken" gorm:"column:physical_damage_taken;type:int"`                         // 承受的物理伤害
	ProfileIcon                    int       `json:"profileIcon" gorm:"column:profile_icon;type:smallint"`                                     // 头像图标ID
	Puuid                          string    `json:"puuid" gorm:"column:puuid;type:varchar(100)"`                                              // 参与者UUID
	QuadraKills                    int       `json:"quadraKills" gorm:"column:quadra_kills;type:smallint"`                                     // 四杀数
	RiotIdName                     string    `json:"riotIdName" gorm:"column:riot_id_name;type:varchar(100)"`                                  // Riot ID名称
	RiotIdTagline                  string    `json:"riotIdTagline" gorm:"column:riot_id_tagline;type:varchar(100)"`                            // Riot ID标签
	Role                           string    `json:"role" gorm:"column:role;type:varchar(100)"`                                                // 角色
	SightWardsBoughtInGame         int       `json:"sightWardsBoughtInGame" gorm:"column:sight_wards_bought_in_game;type:smallint"`            // 购买视野守卫数
	Spell1Casts                    int       `json:"spell1Casts" gorm:"column:spell1_casts;type:smallint"`                                     // 技能1释放次数
	Spell2Casts                    int       `json:"spell2Casts" gorm:"column:spell2_casts;type:smallint"`                                     // 技能2释放次数
	Spell3Casts                    int       `json:"spell3Casts" gorm:"column:spell3_casts;type:smallint"`                                     // 技能3释放次数
	Spell4Casts                    int       `json:"spell4Casts" gorm:"column:spell4_casts;type:smallint"`                                     // 技能4释放次数
	Summoner1Casts                 int       `json:"summoner1Casts" gorm:"column:summoner1_casts;type:smallint"`                               // 召唤师技能1释放次数
	Summoner1Id                    int       `json:"summoner1Id" gorm:"column:summoner1_id;type:smallint"`                                     // 召唤师技能1ID
	Summoner2Casts                 int       `json:"summoner2Casts" gorm:"column:summoner2_casts;type:smallint"`                               // 召唤师技能2释放次数
	Summoner2Id                    int       `json:"summoner2Id" gorm:"column:summoner2_id;type:smallint"`                                     // 召唤师技能2ID
	SummonerId                     string    `json:"summonerId" gorm:"column:summoner_id;index;type:varchar(100)"`                             // 召唤师ID
	SummonerLevel                  int       `json:"summonerLevel" gorm:"column:summoner_level;type:smallint"`                                 // 召唤师等级
	SummonerName                   string    `json:"summonerName" gorm:"column:summoner_name;type:varchar(100)"`                               // 召唤师名称
	TeamEarlySurrendered           bool      `json:"teamEarlySurrendered" gorm:"column:team_early_surrendered"`                                // 队伍提前投降
	TeamId                         int       `json:"teamId" gorm:"column:team_id;type:smallint"`                                               // 队伍ID
	TeamPosition                   string    `json:"teamPosition" gorm:"column:team_position;type:varchar(100)"`                               // 队伍位置
	TimeCCingOthers                int       `json:"timeCCingOthers" gorm:"column:time_ccing_others;type:smallint"`                            // 控制敌方英雄时间
	TimePlayed                     int       `json:"timePlayed" gorm:"column:time_played;type:smallint"`                                       // 游戏时间
	TotalDamageDealt               int       `json:"totalDamageDealt" gorm:"column:total_damage_dealt;type:int"`                               // 总伤害输出
	TotalDamageDealtToChampions    int       `json:"totalDamageDealtToChampions" gorm:"column:total_damage_dealt_to_champions;type:int"`       // 对英雄造成的总伤害
	TotalDamageShieldedOnTeammates int       `json:"totalDamageShieldedOnTeammates" gorm:"column:total_damage_shielded_on_teammates;type:int"` // 对队友护盾总量
	TotalDamageTaken               int       `json:"totalDamageTaken" gorm:"column:total_damage_taken;type:int"`                               // 承受的总伤害
	TotalHeal                      int       `json:"totalHeal" gorm:"column:total_heal;type:int"`                                              // 总治疗量
	TotalHealsOnTeammates          int       `json:"totalHealsOnTeammates" gorm:"column:total_heals_on_teammates;type:int"`                    // 对队友的治疗量
	TotalMinionsKilled             int       `json:"totalMinionsKilled" gorm:"column:total_minions_killed;type:smallint"`                      // 总击杀小兵数
	TotalTimeCCDealt               int       `json:"totalTimeCCDealt" gorm:"column:total_time_cc_dealt;type:smallint"`                         // 总控制时间
	TotalTimeSpentDead             int       `json:"totalTimeSpentDead" gorm:"column:total_time_spent_dead;type:smallint"`                     // 总死亡时间
	TotalUnitsHealed               int       `json:"totalUnitsHealed" gorm:"column:total_units_healed;type:smallint"`                          // 总治疗单位数
	TripleKills                    int       `json:"tripleKills" gorm:"column:triple_kills;type:smallint"`                                     // 三杀数
	TrueDamageDealt                int       `json:"trueDamageDealt" gorm:"column:true_damage_dealt;type:int"`                                 // 真实伤害输出
	TrueDamageDealtToChampions     int       `json:"trueDamageDealtToChampions" gorm:"column:true_damage_dealt_to_champions;type:int"`         // 对英雄造成的真实伤害
	TrueDamageTaken                int       `json:"trueDamageTaken" gorm:"column:true_damage_taken;type:int"`                                 // 承受的真实伤害
	TurretKills                    int       `json:"turretKills" gorm:"column:turret_kills;type:smallint"`                                     // 击杀防御塔数
	TurretTakedowns                int       `json:"turretTakedowns" gorm:"column:turret_takedowns;type:smallint"`                             // 击杀防御塔参与数
	TurretsLost                    int       `json:"turretsLost" gorm:"column:turrets_lost;type:smallint"`                                     // 失去的防御塔数
	UnrealKills                    int       `json:"unrealKills" gorm:"column:unreal_kills;type:smallint"`                                     // 不可思议的击杀数
	VisionScore                    int       `json:"visionScore" gorm:"column:vision_score;type:smallint"`                                     // 视野得分
	VisionWardsBoughtInGame        int       `json:"visionWardsBoughtInGame" gorm:"column:vision_wards_bought_in_game;type:smallint"`          // 购买视野守卫数
	WardsKilled                    int       `json:"wardsKilled" gorm:"column:wards_killed;type:smallint"`                                     // 摧毁守卫数
	WardsPlaced                    int       `json:"wardsPlaced" gorm:"column:wards_placed;type:smallint"`                                     // 放置守卫数
	Win                            bool      `json:"win" gorm:"column:win"`                                                                    // 是否获胜
}


type PerksDto struct {
	StatPerks struct {
		Defense int `json:"defense"` // 防御属性
		Flex    int `json:"flex"`    // 弹性属性
		Offense int `json:"offense"` // 进攻属性
	} `json:"statPerks"`
	
	Styles []struct {
		Description string `json:"description"` // 描述
		Selections  []struct {
			Perk int `json:"perk"` // 符文
			Var1 int `json:"var1"` // 变量1
			Var2 int `json:"var2"` // 变量2
			Var3 int `json:"var3"` // 变量3
		} `json:"selections"`
		Style int `json:"style"` // 符文样式
	} `json:"styles"`
}

func (p *PerksDto) parsePerksMeta(buff *PerksORM) {
	buff.StatDefence = p.StatPerks.Defense
	buff.StatFlex = p.StatPerks.Flex
	buff.StatOffense = p.StatPerks.Offense
	for _, sty := range p.Styles {
		switch sty.Description {
		case "primaryStyle":
			buff.PriStyle = sty.Style
			buff.PriSelection0 = sty.Selections[0].Perk
			buff.PriSelection1 = sty.Selections[1].Perk
			buff.PriSelection2 = sty.Selections[2].Perk
			buff.PriSelection3 = sty.Selections[3].Perk
		case "subStyle":
			buff.SubStyle = sty.Style
			buff.SubSelection0 = sty.Selections[0].Perk
			buff.SubSelection1 = sty.Selections[1].Perk
		}
	}
}

type PerksORM struct {
	StatDefence   int `gorm:"column:stat_defence;type:smallint"`    // 防御属性
	StatFlex      int `gorm:"column:stat_flex;type:smallint"`       // 弹性属性
	StatOffense   int `gorm:"column:stat_offense;type:smallint"`    // 进攻属性
	PriStyle      int `gorm:"column:pri_style;type:smallint"`       // 主要样式
	PriSelection0 int `gorm:"column:pri_selection_0;type:smallint"` // 主要选择0
	PriSelection1 int `gorm:"column:pri_selection_1;type:smallint"` // 主要选择1
	PriSelection2 int `gorm:"column:pri_selection_2;type:smallint"` // 主要选择2
	PriSelection3 int `gorm:"column:pri_selection_3;type:smallint"` // 主要选择3
	SubStyle      int `gorm:"column:sub_style;type:smallint"`       // 次要样式
	SubSelection0 int `gorm:"column:sub_selection_0;type:smallint"` // 次要选择0
	SubSelection1 int `gorm:"column:sub_selection_1;type:smallint"` // 次要选择1
}

