package riotmodel

import (
	"fmt"
)

// Participant won't store in db anymore
type Participant struct {
	Assists                        int         `json:"assists" gorm:"column:assists;type:smallint"`                                              // 助攻数
	BaronKills                     int         `json:"baronKills" gorm:"column:baron_kills;type:smallint"`                                       // 击杀大龙数
	BountyLevel                    int         `json:"bountyLevel" gorm:"column:bounty_level;type:smallint"`                                     // 赏金等级
	Challenges                     *Challenges `json:"challenges" gorm:"embedded;embeddedPrefix:challenges_"`                                    // 达成的挑战数据
	ChampExperience                int         `json:"champExperience" gorm:"column:champ_experience;type:int"`                                  // 英雄经验
	ChampLevel                     int         `json:"champLevel" gorm:"column:champ_level;type:smallint"`                                       // 英雄等级
	ChampionId                     int         `json:"championId" gorm:"column:champion_id;type:smallint"`                                       // 英雄ID
	ChampionName                   string      `json:"championName" gorm:"column:champion_name;type:varchar(100)"`                               // 英雄名称
	ChampionTransform              int         `json:"championTransform" gorm:"column:champion_transform;type:smallint"`                         // 英雄转化（Kayn的转型）
	ConsumablesPurchased           int         `json:"consumablesPurchased" gorm:"column:consumables_purchased;type:smallint"`                   // 购买消耗品数
	DamageDealtToBuildings         int         `json:"damageDealtToBuildings" gorm:"column:damage_dealt_to_buildings;type:int"`                  // 对建筑造成的伤害
	DamageDealtToObjectives        int         `json:"damageDealtToObjectives" gorm:"column:damage_dealt_to_objectives;type:int"`                // 对目标造成的伤害
	DamageDealtToTurrets           int         `json:"damageDealtToTurrets" gorm:"column:damage_dealt_to_turrets;type:int"`                      // 对防御塔造成的伤害
	DamageSelfMitigated            int         `json:"damageSelfMitigated" gorm:"column:damage_self_mitigated;type:int"`                         // 减少的伤害
	Deaths                         int         `json:"deaths" gorm:"column:deaths;type:smallint"`                                                // 死亡数
	DetectorWardsPlaced            int         `json:"detectorWardsPlaced" gorm:"column:detector_wards_placed;type:smallint"`                    // 放置控制视野守卫数
	DoubleKills                    int         `json:"doubleKills" gorm:"column:double_kills;type:smallint"`                                     // 双杀数
	DragonKills                    int         `json:"dragonKills" gorm:"column:dragon_kills;type:smallint"`                                     // 击杀小龙数
	FirstBloodAssist               bool        `json:"firstBloodAssist" gorm:"column:first_blood_assist"`                                        // 是否协助一血
	FirstBloodKill                 bool        `json:"firstBloodKill" gorm:"column:first_blood_kill"`                                            // 是否击杀一血
	FirstTowerAssist               bool        `json:"firstTowerAssist" gorm:"column:first_tower_assist"`                                        // 是否协助一塔
	FirstTowerKill                 bool        `json:"firstTowerKill" gorm:"column:first_tower_kill"`                                            // 是否击杀一塔
	GameEndedInEarlySurrender      bool        `json:"gameEndedInEarlySurrender" gorm:"column:game_ended_in_early_surrender"`                    // 是否提前投降结束游戏
	GameEndedInSurrender           bool        `json:"gameEndedInSurrender" gorm:"column:game_ended_in_surrender"`                               // 是否投降结束游戏
	GoldEarned                     int         `json:"goldEarned" gorm:"column:gold_earned;type:int"`                                            // 获得的金币
	GoldSpent                      int         `json:"goldSpent" gorm:"column:gold_spent;type:int"`                                              // 花费的金币
	IndividualPosition             string      `json:"individualPosition" gorm:"column:individual_position;type:varchar(100)"`                   // 个人位置
	InhibitorKills                 int         `json:"inhibitorKills" gorm:"column:inhibitor_kills;type:smallint"`                               // 击杀抑制塔数
	InhibitorTakedowns             int         `json:"inhibitorTakedowns" gorm:"column:inhibitor_takedowns;type:smallint"`                       // 击杀抑制塔参与数
	InhibitorsLost                 int         `json:"inhibitorsLost" gorm:"column:inhibitors_lost;type:smallint"`                               // 失去的抑制塔数
	Item0                          int         `json:"item0" gorm:"column:item0;type:int"`                                                       // 物品0
	Item1                          int         `json:"item1" gorm:"column:item1;type:int"`                                                       // 物品1
	Item2                          int         `json:"item2" gorm:"column:item2;type:int"`                                                       // 物品2
	Item3                          int         `json:"item3" gorm:"column:item3;type:int"`                                                       // 物品3
	Item4                          int         `json:"item4" gorm:"column:item4;type:int"`                                                       // 物品4
	Item5                          int         `json:"item5" gorm:"column:item5;type:int"`                                                       // 物品5
	Item6                          int         `json:"item6" gorm:"column:item6;type:int"`                                                       // 物品6
	ItemsPurchased                 int         `json:"itemsPurchased" gorm:"column:items_purchased;type:smallint"`                               // 购买物品数
	KillingSprees                  int         `json:"killingSprees" gorm:"column:killing_sprees;type:smallint"`                                 // 连杀数
	Kills                          int         `json:"kills" gorm:"column:kills;type:smallint"`                                                  // 击杀数
	Lane                           string      `json:"lane" gorm:"column:lane;type:varchar(100)"`                                                // 路线
	LargestCriticalStrike          int         `json:"largestCriticalStrike" gorm:"column:largest_critical_strike;type:smallint"`                // 最大暴击伤害
	LargestKillingSpree            int         `json:"largestKillingSpree" gorm:"column:largest_killing_spree;type:smallint"`                    // 最大连杀数
	LargestMultiKill               int         `json:"largestMultiKill" gorm:"column:largest_multi_kill;type:smallint"`                          // 最大多杀数
	LongestTimeSpentLiving         int         `json:"longestTimeSpentLiving" gorm:"column:longest_time_spent_living;type:smallint"`             // 存活最长时间
	MagicDamageDealt               int         `json:"magicDamageDealt" gorm:"column:magic_damage_dealt;type:int"`                               // 魔法伤害输出
	MagicDamageDealtToChampions    int         `json:"magicDamageDealtToChampions" gorm:"column:magic_damage_dealt_to_champions;type:int"`       // 对英雄造成的魔法伤害
	MagicDamageTaken               int         `json:"magicDamageTaken" gorm:"column:magic_damage_taken;type:int"`                               // 承受的魔法伤害
	NeutralMinionsKilled           int         `json:"neutralMinionsKilled" gorm:"column:neutral_minions_killed;type:smallint"`                  // 击杀中立生物数
	NexusKills                     int         `json:"nexusKills" gorm:"column:nexus_kills;type:smallint"`                                       // 击杀水晶数
	NexusTakedowns                 int         `json:"nexusTakedowns" gorm:"column:nexus_takedowns;type:smallint"`                               // 击杀水晶参与数
	NexusLost                      int         `json:"nexusLost" gorm:"column:nexus_lost;type:smallint"`                                         // 失去的水晶数
	ObjectivesStolen               int         `json:"objectivesStolen" gorm:"column:objectives_stolen;type:smallint"`                           // 偷取目标数
	ObjectivesStolenAssists        int         `json:"objectivesStolenAssists" gorm:"column:objectives_stolen_assists;type:smallint"`            // 偷取目标参与数
	ParticipantId                  int         `json:"participantId" gorm:"column:participant_id;type:smallint"`                                 // 参与者ID
	PentaKills                     int         `json:"pentaKills" gorm:"column:penta_kills;type:smallint"`                                       // 五杀数
	PerksMeta                      *PerksDto   `json:"perks" gorm:"-"`                                                                           // 符文
	PerksSTR                       string      `gorm:"column:perks"`                                                                             // 符文
	PhysicalDamageDealt            int         `json:"physicalDamageDealt" gorm:"column:physical_damage_dealt;type:int"`                         // 物理伤害输出
	PhysicalDamageDealtToChampions int         `json:"physicalDamageDealtToChampions" gorm:"column:physical_damage_dealt_to_champions;type:int"` // 对英雄造成的物理伤害
	PhysicalDamageTaken            int         `json:"physicalDamageTaken" gorm:"column:physical_damage_taken;type:int"`                         // 承受的物理伤害
	ProfileIcon                    int         `json:"profileIcon" gorm:"column:profile_icon;type:int"`                                          // 头像图标ID
	Puuid                          string      `json:"puuid" gorm:"column:puuid;type:varchar(100)"`                                              // 参与者UUID
	QuadraKills                    int         `json:"quadraKills" gorm:"column:quadra_kills;type:smallint"`                                     // 四杀数
	RiotIdName                     string      `json:"riotIdName" gorm:"column:riot_id_name;type:varchar(100)"`                                  // Riot ID名称
	RiotIdTagline                  string      `json:"riotIdTagline" gorm:"column:riot_id_tagline;type:varchar(100)"`                            // Riot ID标签
	Role                           string      `json:"role" gorm:"column:role;type:varchar(100)"`                                                // 角色
	SightWardsBoughtInGame         int         `json:"sightWardsBoughtInGame" gorm:"column:sight_wards_bought_in_game;type:smallint"`            // 购买视野守卫数
	Spell1Casts                    int         `json:"spell1Casts" gorm:"column:spell1_casts;type:smallint"`                                     // 技能1释放次数
	Spell2Casts                    int         `json:"spell2Casts" gorm:"column:spell2_casts;type:smallint"`                                     // 技能2释放次数
	Spell3Casts                    int         `json:"spell3Casts" gorm:"column:spell3_casts;type:smallint"`                                     // 技能3释放次数
	Spell4Casts                    int         `json:"spell4Casts" gorm:"column:spell4_casts;type:smallint"`                                     // 技能4释放次数
	Summoner1Casts                 int         `json:"summoner1Casts" gorm:"column:summoner1_casts;type:smallint"`                               // 召唤师技能1释放次数
	Summoner1Id                    int         `json:"summoner1Id" gorm:"column:summoner1_id;type:smallint"`                                     // 召唤师技能1ID
	Summoner2Casts                 int         `json:"summoner2Casts" gorm:"column:summoner2_casts;type:smallint"`                               // 召唤师技能2释放次数
	Summoner2Id                    int         `json:"summoner2Id" gorm:"column:summoner2_id;type:smallint"`                                     // 召唤师技能2ID
	SummonerId                     string      `json:"summonerId" gorm:"column:summoner_id;index;type:varchar(100)"`                             // 召唤师ID
	SummonerLevel                  int         `json:"summonerLevel" gorm:"column:summoner_level;type:smallint"`                                 // 召唤师等级
	SummonerName                   string      `json:"summonerName" gorm:"column:summoner_name;type:varchar(100)"`                               // 召唤师名称
	TeamEarlySurrendered           bool        `json:"teamEarlySurrendered" gorm:"column:team_early_surrendered"`                                // 队伍提前投降
	TeamId                         int         `json:"teamId" gorm:"column:team_id;type:smallint"`                                               // 队伍ID
	TeamPosition                   string      `json:"teamPosition" gorm:"column:team_position;type:varchar(100)"`                               // 队伍位置
	TimeCCingOthers                int         `json:"timeCCingOthers" gorm:"column:time_ccing_others;type:smallint"`                            // 控制敌方英雄时间
	TimePlayed                     int         `json:"timePlayed" gorm:"column:time_played;type:smallint"`                                       // 游戏时间
	TotalDamageDealt               int         `json:"totalDamageDealt" gorm:"column:total_damage_dealt;type:int"`                               // 总伤害输出
	TotalDamageDealtToChampions    int         `json:"totalDamageDealtToChampions" gorm:"column:total_damage_dealt_to_champions;type:int"`       // 对英雄造成的总伤害
	TotalDamageShieldedOnTeammates int         `json:"totalDamageShieldedOnTeammates" gorm:"column:total_damage_shielded_on_teammates;type:int"` // 对队友护盾总量
	TotalDamageTaken               int         `json:"totalDamageTaken" gorm:"column:total_damage_taken;type:int"`                               // 承受的总伤害
	TotalHeal                      int         `json:"totalHeal" gorm:"column:total_heal;type:int"`                                              // 总治疗量
	TotalHealsOnTeammates          int         `json:"totalHealsOnTeammates" gorm:"column:total_heals_on_teammates;type:int"`                    // 对队友的治疗量
	TotalMinionsKilled             int         `json:"totalMinionsKilled" gorm:"column:total_minions_killed;type:smallint"`                      // 补刀总数
	TotalTimeCCDealt               int         `json:"totalTimeCCDealt" gorm:"column:total_time_cc_dealt;type:smallint"`                         // 总控制时间
	TotalTimeSpentDead             int         `json:"totalTimeSpentDead" gorm:"column:total_time_spent_dead;type:smallint"`                     // 总死亡时间
	TotalUnitsHealed               int         `json:"totalUnitsHealed" gorm:"column:total_units_healed;type:smallint"`                          // 总治疗单位数
	TripleKills                    int         `json:"tripleKills" gorm:"column:triple_kills;type:smallint"`                                     // 三杀数
	TrueDamageDealt                int         `json:"trueDamageDealt" gorm:"column:true_damage_dealt;type:int"`                                 // 真实伤害输出
	TrueDamageDealtToChampions     int         `json:"trueDamageDealtToChampions" gorm:"column:true_damage_dealt_to_champions;type:int"`         // 对英雄造成的真实伤害
	TrueDamageTaken                int         `json:"trueDamageTaken" gorm:"column:true_damage_taken;type:int"`                                 // 承受的真实伤害
	TurretKills                    int         `json:"turretKills" gorm:"column:turret_kills;type:smallint"`                                     // 击杀防御塔数
	TurretTakedowns                int         `json:"turretTakedowns" gorm:"column:turret_takedowns;type:smallint"`                             // 击杀防御塔参与数
	TurretsLost                    int         `json:"turretsLost" gorm:"column:turrets_lost;type:smallint"`                                     // 失去的防御塔数
	UnrealKills                    int         `json:"unrealKills" gorm:"column:unreal_kills;type:smallint"`                                     // 不可思议的击杀数
	VisionScore                    int         `json:"visionScore" gorm:"column:vision_score;type:smallint"`                                     // 视野得分
	VisionWardsBoughtInGame        int         `json:"visionWardsBoughtInGame" gorm:"column:vision_wards_bought_in_game;type:smallint"`          // 购买视野守卫数
	WardsKilled                    int         `json:"wardsKilled" gorm:"column:wards_killed;type:smallint"`                                     // 摧毁守卫数
	WardsPlaced                    int         `json:"wardsPlaced" gorm:"column:wards_placed;type:smallint"`                                     // 放置守卫数
	Win                            bool        `json:"win" gorm:"column:win"`                                                                    // 是否获胜
}

// Challenges (not exist in CHERRY mode)
type Challenges struct {
	AssistStreakCount                        uint    `json:"12AssistStreakCount" gorm:"column:assist_streak_count;type:smallint"`                                                   // 12次助攻连续次数
	AbilityUses                              uint    `json:"abilityUses" gorm:"column:ability_uses;type:smallint"`                                                                  // 使用技能次数
	AcesBefore15Minutes                      uint    `json:"acesBefore15Minutes" gorm:"column:aces_before_15_minutes;type:smallint"`                                                // 15分钟前完成团灭次数
	AlliedJungleMonsterKills                 float32 `json:"alliedJungleMonsterKills" gorm:"column:allied_jungle_monster_kills;type:float"`                                         // 友方野怪击杀数
	BaronTakedowns                           uint    `json:"baronTakedowns" gorm:"column:baron_takedowns;type:smallint"`                                                            // 击杀男爵次数
	BlastConeOppositeOpponentCount           uint    `json:"blastConeOppositeOpponentCount" gorm:"column:blast_cone_opposite_opponent_count;type:smallint"`                         // 使用位移锥击飞对方次数
	BountyGold                               uint    `json:"bountyGold" gorm:"column:bounty_gold;type:smallint"`                                                                    // 赏金金币总数
	BuffsStolen                              uint    `json:"buffsStolen" gorm:"column:buffs_stolen;type:smallint"`                                                                  // 偷取增益效果次数
	CompleteSupportQuestInTime               uint    `json:"completeSupportQuestInTime" gorm:"column:complete_support_quest_in_time;type:smallint"`                                 // 在时间内完成辅助任务次数
	ControlWardsPlaced                       uint    `json:"controlWardsPlaced" gorm:"column:control_wards_placed;type:smallint"`                                                   // 放置控制守卫次数
	DamagePerMinute                          float32 `json:"damagePerMinute" gorm:"column:damage_per_minute;type:float"`                                                            // 每分钟造成伤害
	DamageTakenOnTeamPercentage              float32 `json:"damageTakenOnTeamPercentage" gorm:"column:damage_taken_on_team_percentage;type:float"`                                  // 团队总伤害承受占比
	DancedWithRiftHerald                     uint    `json:"dancedWithRiftHerald" gorm:"column:danced_with_rift_herald;type:smallint"`                                              // 与峡谷先锋共舞次数
	DeathsByEnemyChamps                      uint    `json:"deathsByEnemyChamps" gorm:"column:deaths_by_enemy_champs;type:smallint"`                                                // 被敌方英雄击杀次数
	DodgeSkillShotsSmallWindow               uint    `json:"dodgeSkillShotsSmallWindow" gorm:"column:dodge_skill_shots_small_window;type:smallint"`                                 // 在短时间内躲避技能弹道次数
	DoubleAces                               uint    `json:"doubleAces" gorm:"column:double_aces;type:smallint"`                                                                    // 双团灭次数
	DragonTakedowns                          uint    `json:"dragonTakedowns" gorm:"column:dragon_takedowns;type:smallint"`                                                          // 击杀龙次数
	EarlyLaningPhaseGoldExpAdvantage         uint    `json:"earlyLaningPhaseGoldExpAdvantage" gorm:"column:early_laning_phase_gold_exp_advantage;type:smallint"`                    // 早期金币和经验优势
	EffectiveHealAndShielding                float32 `json:"effectiveHealAndShielding" gorm:"column:effective_heal_and_shielding;type:float"`                                       // 有效治疗和护盾总量
	ElderDragonKillsWithOpposingSoul         uint    `json:"elderDragonKillsWithOpposingSoul" gorm:"column:elder_dragon_kills_with_opposing_soul;type:smallint"`                    // 远古龙斩杀次数?
	ElderDragonMultikills                    uint    `json:"elderDragonMultikills" gorm:"column:elder_dragon_multikills;type:smallint"`                                             // 远古龙多杀次数
	EnemyChampionImmobilizations             uint    `json:"enemyChampionImmobilizations" gorm:"column:enemy_champion_immobilizations;type:smallint"`                               // 对敌方英雄的控制效果次数
	EnemyJungleMonsterKills                  float32 `json:"enemyJungleMonsterKills" gorm:"column:enemy_jungle_monster_kills;type:float"`                                           // 偷野次数
	EpicMonsterKillsNearEnemyJungler         uint    `json:"epicMonsterKillsNearEnemyJungler" gorm:"column:epic_monster_kills_near_enemy_jungler;type:smallint"`                    // 在敌方丛林附近击杀史诗怪物次数
	EpicMonsterKillsWithin30SecondsOfSpawn   uint    `json:"epicMonsterKillsWithin30SecondsOfSpawn" gorm:"column:epic_monster_kills_within_30_seconds_of_spawn;type:smallint"`      // 在史诗怪物生成后30秒内击杀次数
	EpicMonsterSteals                        uint    `json:"epicMonsterSteals" gorm:"column:epic_monster_steals;type:smallint"`                                                     // 偷取史诗怪物次数
	EpicMonsterStolenWithoutSmite            uint    `json:"epicMonsterStolenWithoutSmite" gorm:"column:epic_monster_stolen_without_smite;type:smallint"`                           // 不使用惩戒技能偷取史诗怪物次数
	FirstTurretKilled                        float32 `json:"firstTurretKilled" gorm:"column:first_turret_killed;type:float"`                                                        // 一血塔摧毁标记 || 时间
	FirstTurretKilledTime                    float32 `json:"firstTurretKilledTime" gorm:"column:first_turret_killed_time;type:float"`                                               // 一血塔摧毁时间
	FlawlessAces                             uint    `json:"flawlessAces" gorm:"column:flawless_aces;type:smallint"`                                                                // 完美的团灭次数
	FullTeamTakedown                         uint    `json:"fullTeamTakedown" gorm:"column:full_team_takedown;type:smallint"`                                                       // 整个团队击杀次数
	GameLength                               float32 `json:"gameLength" gorm:"column:game_length;type:float"`                                                                       // 比赛时长
	GetTakedownsInAllLanesEarlyJungleAsLaner uint    `json:"getTakedownsInAllLanesEarlyJungleAsLaner" gorm:"column:get_takedowns_in_all_lanes_early_jungle_as_laner;type:smallint"` // 早期在所有路线和丛林击杀次数
	GoldPerMinute                            float32 `json:"goldPerMinute" gorm:"column:gold_per_minute;type:float"`                                                                // 每分钟获得金币
	HadOpenNexus                             uint    `json:"hadOpenNexus" gorm:"column:had_open_nexus;type:smallint"`                                                               // 已击破敌方水晶次数
	ImmobilizeAndKillWithAlly                uint    `json:"immobilizeAndKillWithAlly" gorm:"column:immobilize_and_kill_with_ally;type:smallint"`                                   // 与队友定身击杀次数
	InitialBuffCount                         uint    `json:"initialBuffCount" gorm:"column:initial_buff_count;type:smallint"`                                                       // 初始增益效果次数
	InitialCrabCount                         uint    `json:"initialCrabCount" gorm:"column:initial_crab_count;type:smallint"`                                                       // 初始螃蟹击杀次数
	JungleCsBefore10Minutes                  float32 `json:"jungleCsBefore10Minutes" gorm:"column:jungle_cs_before_10_minutes;type:float"`                                          // 10分钟前击杀丛林怪物次数
	JunglerTakedownsNearDamagedEpicMonster   uint    `json:"junglerTakedownsNearDamagedEpicMonster" gorm:"column:jungler_takedowns_near_damaged_epic_monster;type:smallint"`        // 击杀受损史诗怪物附近的丛林怪物次数
	KTurretsDestroyedBeforePlatesFall        uint    `json:"kTurretsDestroyedBeforePlatesFall" gorm:"column:k_turrets_destroyed_before_plates_fall;type:smallint"`                  // 在防御塔血条降低前击杀塔次数
	KDA                                      float32 `json:"kda" gorm:"column:kda;type:float"`                                                                                      // KDA比分
	KillAfterHiddenWithAlly                  uint    `json:"killAfterHiddenWithAlly" gorm:"column:kill_after_hidden_with_ally;type:smallint"`                                       // 与队友隐藏后击杀次数
	KillParticipation                        float32 `json:"killParticipation" gorm:"column:kill_participation;type:float"`                                                         // 击杀参与率
	KilledChampTookFullTeamDamageSurvived    uint    `json:"killedChampTookFullTeamDamageSurvived" gorm:"column:killed_champ_took_full_team_damage_survived;type:smallint"`         // 击杀敌方英雄并全队伤害幸存次数
	KillingSprees                            uint    `json:"killingSprees" gorm:"column:killing_sprees;type:smallint"`                                                              // 连续击杀次数
	KillsNearEnemyTurret                     uint    `json:"killsNearEnemyTurret" gorm:"column:kills_near_enemy_turret;type:smallint"`                                              // 敌方防御塔附近击杀次数
	KillsOnOtherLanesEarlyJungleAsLaner      uint    `json:"killsOnOtherLanesEarlyJungleAsLaner" gorm:"column:kills_on_other_lanes_early_jungle_as_laner;type:smallint"`            // 早期在其他线路和丛林击杀次数
	KillsOnRecentlyHealedByAramPack          uint    `json:"killsOnRecentlyHealedByAramPack" gorm:"column:kills_on_recently_healed_by_aram_pack;type:smallint"`                     // 击杀刚吃完ARAM血包次数
	KillsUnderOwnTurret                      uint    `json:"killsUnderOwnTurret" gorm:"column:kills_under_own_turret;type:smallint"`                                                // 己方防御塔下击杀次数
	KillsWithHelpFromEpicMonster             uint    `json:"killsWithHelpFromEpicMonster" gorm:"column:kills_with_help_from_epic_monster;type:smallint"`                            // 受到史诗怪物助攻击杀次数
	KnockEnemyIntoTeamAndKill                uint    `json:"knockEnemyIntoTeamAndKill" gorm:"column:knock_enemy_into_team_and_kill;type:smallint"`                                  // 击退敌方英雄并与全队合作击杀次数
	LandSkillShotsEarlyGame                  uint    `json:"landSkillShotsEarlyGame" gorm:"column:land_skill_shots_early_game;type:smallint"`                                       // 早期游戏中命中技能次数
	LaneMinionsFirst10Minutes                uint    `json:"laneMinionsFirst10Minutes" gorm:"column:lane_minions_first_10_minutes;type:smallint"`                                   // 前10分钟击杀兵线小兵次数
	LaningPhaseGoldExpAdvantage              uint    `json:"laningPhaseGoldExpAdvantage" gorm:"column:laning_phase_gold_exp_advantage;type:smallint"`                               // 上路金币和经验优势
	LegendaryCount                           uint    `json:"legendaryCount" gorm:"column:legendary_count;type:smallint"`                                                            // 超神次数
	LostAnInhibitor                          uint    `json:"lostAnInhibitor" gorm:"column:lost_an_inhibitor;type:smallint"`                                                         // 丢失兵营次数
	MaxCsAdvantageOnLaneOpponent             float32 `json:"maxCsAdvantageOnLaneOpponent" gorm:"column:max_cs_advantage_on_lane_opponent;type:float"`                               // 对线补刀差
	MaxKillDeficit                           uint    `json:"maxKillDeficit" gorm:"column:max_kill_deficit;type:smallint"`                                                           // 最大击杀差距次数
	MaxLevelLeadLaneOpponent                 uint    `json:"maxLevelLeadLaneOpponent" gorm:"column:max_level_lead_lane_opponent;type:smallint"`                                     // 对线等级差
	MejaisFullStackInTime                    uint    `json:"mejaisFullStackInTime" gorm:"column:mejais_full_stack_in_time;type:smallint"`                                           // 在时间内堆满杀人书的次数
	MoreEnemyJungleThanOpponent              float32 `json:"moreEnemyJungleThanOpponent" gorm:"column:more_enemy_jungle_than_opponent;type:float"`                                  // 击杀更多敌方丛林怪物次数
	MultiKillOneSpell                        uint    `json:"multiKillOneSpell" gorm:"column:multi_kill_one_spell;type:smallint"`                                                    // 单技能多杀次数
	MultiTurretRiftHeraldCount               uint    `json:"multiTurretRiftHeraldCount" gorm:"column:multi_turret_rift_herald_count;type:smallint"`                                 // 峡谷先锋击杀多个防御塔次数
	Multikills                               uint    `json:"multikills" gorm:"column:multikills;type:smallint"`                                                                     // 多杀次数
	MultikillsAfterAggressiveFlash           uint    `json:"multikillsAfterAggressiveFlash" gorm:"column:multikills_after_aggressive_flash;type:smallint"`                          // 在激进闪现后连续击杀多名敌方英雄次数
	MythicItemUsed                           uint    `json:"mythicItemUsed" gorm:"column:mythic_item_used;type:smallint"`                                                           // 使用神话物品次数
	OuterTurretExecutesBefore10Minutes       uint    `json:"outerTurretExecutesBefore10Minutes" gorm:"column:outer_turret_executes_before_10_minutes;type:smallint"`                // 前10分钟击杀外塔次数
	OutnumberedKills                         uint    `json:"outnumberedKills" gorm:"column:outnumbered_kills;type:smallint"`                                                        // 人数劣势击杀次数
	OutnumberedNexusKill                     uint    `json:"outnumberedNexusKill" gorm:"column:outnumbered_nexus_kill;type:smallint"`                                               // 人数劣势击杀水晶次数
	PerfectDragonSoulsTaken                  uint    `json:"perfectDragonSoulsTaken" gorm:"column:perfect_dragon_souls_taken;type:smallint"`                                        // 纯种龙魂获取次数
	PerfectGame                              uint    `json:"perfectGame" gorm:"column:perfect_game;type:smallint"`                                                                  // 完美比赛次数
	PickKillWithAlly                         uint    `json:"pickKillWithAlly" gorm:"column:pick_kill_with_ally;type:smallint"`                                                      // 与队友合作击杀次数
	PoroExplosions                           uint    `json:"poroExplosions" gorm:"column:poro_explosions;type:smallint"`                                                            // 魄罗爆炸次数
	QuickCleanse                             uint    `json:"quickCleanse" gorm:"column:quick_cleanse;type:smallint"`                                                                // 快速解控次数
	QuickFirstTurret                         uint    `json:"quickFirstTurret" gorm:"column:quick_first_turret;type:smallint"`                                                       // 快速击杀第一座防御塔次数
	QuickSoloKills                           uint    `json:"quickSoloKills" gorm:"column:quick_solo_kills;type:smallint"`                                                           // 快速单杀次数
	RiftHeraldTakedowns                      uint    `json:"riftHeraldTakedowns" gorm:"column:rift_herald_takedowns;type:smallint"`                                                 // 击杀峡谷先锋次数
	SaveAllyFromDeath                        uint    `json:"saveAllyFromDeath" gorm:"column:save_ally_from_death;type:smallint"`                                                    // 拯救队友于死亡边缘次数
	ScuttleCrabKills                         uint    `json:"scuttleCrabKills" gorm:"column:scuttle_crab_kills;type:smallint"`                                                       // 击杀河蟹次数
	ShortestTimeToAceFromFirstTakedown       float32 `json:"shortestTimeToAceFromFirstTakedown" gorm:"column:shortest_time_to_ace_from_first_takedown;type:float"`                  // 从首次击杀到全灭的最短时间
	SkillshotsDodged                         uint    `json:"skillshotsDodged" gorm:"column:skillshots_dodged;type:smallint"`                                                        // 躲避技能弹道次数
	SkillshotsHit                            uint    `json:"skillshotsHit" gorm:"column:skillshots_hit;type:smallint"`                                                              // 命中技能弹道次数
	SnowballsHit                             uint    `json:"snowballsHit" gorm:"column:snowballs_hit;type:smallint"`                                                                // 命中雪球次数
	SoloBaronKills                           uint    `json:"soloBaronKills" gorm:"column:solo_baron_kills;type:smallint"`                                                           // 单人击杀男爵次数
	SoloKills                                uint    `json:"soloKills" gorm:"column:solo_kills;type:smallint"`                                                                      // 单人击杀次数
	StealthWardsPlaced                       uint    `json:"stealthWardsPlaced" gorm:"column:stealth_wards_placed;type:smallint"`                                                   // 放置隐形守卫次数
	SurvivedSingleDigitHpCount               uint    `json:"survivedSingleDigitHpCount" gorm:"column:survived_single_digit_hp_count;type:smallint"`                                 // 幸存单个数字血量次数
	SurvivedThreeImmobilizesInFight          uint    `json:"survivedThreeImmobilizesInFight" gorm:"column:survived_three_immobilizes_in_fight;type:smallint"`                       // 在战斗中幸存三次定身次数
	TakedownOnFirstTurret                    uint    `json:"takedownOnFirstTurret" gorm:"column:takedown_on_first_turret;type:smallint"`                                            // 击杀第一座防御塔次数
	Takedowns                                uint    `json:"takedowns" gorm:"column:takedowns;type:smallint"`                                                                       // 击杀次数
	TakedownsAfterGainingLevelAdvantage      uint    `json:"takedownsAfterGainingLevelAdvantage" gorm:"column:takedowns_after_gaining_level_advantage;type:smallint"`               // 获得等级优势后击杀次数
	TakedownsBeforeJungleMinionSpawn         uint    `json:"takedownsBeforeJungleMinionSpawn" gorm:"column:takedowns_before_jungle_minion_spawn;type:smallint"`                     // 丛林小兵生成前击杀次数
	TakedownsFirstXMinutes                   uint    `json:"takedownsFirstXMinutes" gorm:"column:takedowns_first_x_minutes;type:smallint"`                                          // 前 X 分钟击杀次数
	TakedownsInAlcove                        uint    `json:"takedownsInAlcove" gorm:"column:takedowns_in_alcove;type:smallint"`                                                     // 位于凹角处击杀次数
	TakedownsInEnemyFountain                 uint    `json:"takedownsInEnemyFountain" gorm:"column:takedowns_in_enemy_fountain;type:smallint"`                                      // 位于敌方泉水击杀次数
	TeamBaronKills                           uint    `json:"teamBaronKills" gorm:"column:team_baron_kills;type:smallint"`                                                           // 团队击杀男爵次数
	TeamDamagePercentage                     float32 `json:"teamDamagePercentage" gorm:"column:team_damage_percentage;type:float"`                                                  // 团队伤害占比
	TeamElderDragonKills                     uint    `json:"teamElderDragonKills" gorm:"column:team_elder_dragon_kills;type:smallint"`                                              // 团队击杀古龙次数
	TeamRiftHeraldKills                      uint    `json:"teamRiftHeraldKills" gorm:"column:team_rift_herald_kills;type:smallint"`                                                // 团队击杀先锋女妖次数
	TookLargeDamageSurvived                  uint    `json:"tookLargeDamageSurvived" gorm:"column:took_large_damage_survived;type:smallint"`                                        // 承受大量伤害后幸存次数
	TurretPlatesTaken                        uint    `json:"turretPlatesTaken" gorm:"column:turret_plates_taken;type:smallint"`                                                     // 获取塔皮数
	TurretTakedowns                          uint    `json:"turretTakedowns" gorm:"column:turret_takedowns;type:smallint"`                                                          // 击杀防御塔次数
	TurretsTakenWithRiftHerald               uint    `json:"turretsTakenWithRiftHerald" gorm:"column:turrets_taken_with_rift_herald;type:smallint"`                                 // 峡谷先锋帮助拆除防御塔次数
	TwentyMinionsIn3SecondsCount             uint    `json:"twentyMinionsIn3SecondsCount" gorm:"column:twenty_minions_in_3_seconds_count;type:smallint"`                            // 3秒内击杀20个小兵次数
	TwoWardsOneSweeperCount                  uint    `json:"twoWardsOneSweeperCount" gorm:"column:two_wards_one_sweeper_count;type:smallint"`                                       // 放置两个守卫和一个扫描守卫次数
	UnseenRecalls                            uint    `json:"unseenRecalls" gorm:"column:unseen_recalls;type:smallint"`                                                              // 隐身回城次数
	VisionScoreAdvantageLaneOpponent         float32 `json:"visionScoreAdvantageLaneOpponent" gorm:"column:vision_score_advantage_lane_opponent;type:float"`                        // 与对线对手的视野得分差
	VisionScorePerMinute                     float32 `json:"visionScorePerMinute" gorm:"column:vision_score_per_minute;type:float"`                                                 // 每分钟视野得分
	WardTakedowns                            uint    `json:"wardTakedowns" gorm:"column:ward_takedowns;type:smallint"`                                                              // 摧毁守卫次数
	WardTakedownsBefore20M                   uint    `json:"wardTakedownsBefore20M" gorm:"column:ward_takedowns_before_20m;type:smallint"`                                          // 20分钟前击杀守卫次数
	WardsGuarded                             uint    `json:"wardsGuarded" gorm:"column:wards_guarded;type:smallint"`                                                                // 守卫保护次数
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

func (p *PerksDto) parsePerksMeta() string {
	perk := make(map[string]int)
	for _, sty := range p.Styles {
		switch sty.Description {
		case "primaryStyle":
			perk["priSty"] = sty.Style
			perk["pri0"] = sty.Selections[0].Perk
			perk["pri1"] = sty.Selections[1].Perk
			perk["pri2"] = sty.Selections[2].Perk
			perk["pri3"] = sty.Selections[3].Perk
		case "subStyle":
			perk["subSty"] = sty.Style
			perk["sub0"], perk["sub1"] = sortPerk(sty.Selections[0].Perk, sty.Selections[1].Perk)
		}
	}
	perk["statOffense"] = p.StatPerks.Offense
	perk["statFlex"], perk["statDefence"] = sortPerk(p.StatPerks.Flex, p.StatPerks.Defense)
	
	return fmt.Sprintf("pri:%d,%d,%d,%d,%d sub:%d,%d,%d stat:%d,%d,%d",
		perk["priSty"], perk["pri0"], perk["pri1"], perk["pri2"], perk["pri3"],
		perk["subSty"], perk["sub0"], perk["sub1"],
		perk["statOffense"], perk["statFlex"], perk["statDefence"],
	)
}

func sortPerk(a, b int) (int, int) {
	if a > b {
		return a, b
	}
	return b, a
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
