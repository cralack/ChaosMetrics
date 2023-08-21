package riotmodel

type MatchTimelineDTO struct {
	Metadata MetadataDto `json:"metadata"`
	Info     struct {
		FrameInterval int      `json:"frameInterval"`
		Frames        []*Frame `json:"frames"`
	}
}

type Frame struct {
	Events            []Event                     `json:"events"`
	ParticipantFrames map[string]ParticipantFrame `json:"participantFrames"`
	Timestamp         int                         `json:"timestamp"`
}

type ParticipantFrames map[string]ParticipantFrame

type ParticipantFrame struct {
	ChampionStats            *TLChampionStats `json:"championStats"`            // 冠军统计信息
	CurrentGold              int              `json:"currentGold"`              // 当前金币
	DamageStats              *TLDamageStats   `json:"damageStats"`              // 伤害统计信息
	GoldPerSecond            int              `json:"goldPerSecond"`            // 每秒金币获取量
	JungleMinionsKilled      int              `json:"jungleMinionsKilled"`      // 击杀野怪数量
	Level                    int              `json:"level"`                    // 等级
	MinionsKilled            int              `json:"minionsKilled"`            // 击杀小兵数量
	ParticipantId            int              `json:"participantId"`            // 参与者 ID
	Position                 *Position        `json:"position"`                 // 位置信息
	TimeEnemySpentControlled int              `json:"timeEnemySpentControlled"` // 控制敌方的时间
	TotalGold                int              `json:"totalGold"`                // 总金币
	Xp                       int              `json:"xp"`                       // 经验值
}

type TLChampionStats struct {
	AbilityHaste         int `json:"abilityHaste"`         // 技能急速
	AbilityPower         int `json:"abilityPower"`         // 法术强度
	Armor                int `json:"armor"`                // 护甲
	ArmorPen             int `json:"armorPen"`             // 护甲穿透
	ArmorPenPercent      int `json:"armorPenPercent"`      // 护甲穿透百分比
	AttackDamage         int `json:"attackDamage"`         // 攻击力
	AttackSpeed          int `json:"attackSpeed"`          // 攻击速度
	BonusArmorPenPercent int `json:"bonusArmorPenPercent"` // 附加护甲穿透百分比
	BonusMagicPenPercent int `json:"bonusMagicPenPercent"` // 附加魔法穿透百分比
	CCReduction          int `json:"ccReduction"`          // 韧性
	CooldownReduction    int `json:"cooldownReduction"`    // 冷却时间缩减
	Health               int `json:"health"`               // 生命值
	HealthMax            int `json:"healthMax"`            // 最大生命值
	HealthRegen          int `json:"healthRegen"`          // 生命值回复
	Lifesteal            int `json:"lifesteal"`            // 生命偷取
	MagicPen             int `json:"magicPen"`             // 魔法穿透
	MagicPenPercent      int `json:"magicPenPercent"`      // 魔法穿透百分比
	MagicResist          int `json:"magicResist"`          // 魔法抗性
	MovementSpeed        int `json:"movementSpeed"`        // 移动速度
	Omnivamp             int `json:"omnivamp"`             // 全能吸血
	PhysicalVamp         int `json:"physicalVamp"`         // 物理吸血
	Power                int `json:"power"`                // 能力值
	PowerMax             int `json:"powerMax"`             // 最大能力值
	PowerRegen           int `json:"powerRegen"`           // 能力值回复
	SpellVamp            int `json:"spellVamp"`            // 法术吸血
}

type TLDamageStats struct {
	MagicDamageDone               int `json:"magicDamageDone"`               // 魔法伤害输出
	MagicDamageDoneToChampions    int `json:"magicDamageDoneToChampions"`    // 对英雄的魔法伤害输出
	MagicDamageTaken              int `json:"magicDamageTaken"`              // 承受的魔法伤害
	PhysicalDamageDone            int `json:"physicalDamageDone"`            // 物理伤害输出
	PhysicalDamageDoneToChampions int `json:"physicalDamageDoneToChampions"` // 对英雄的物理伤害输出
	PhysicalDamageTaken           int `json:"physicalDamageTaken"`           // 承受的物理伤害
	TotalDamageDone               int `json:"totalDamageDone"`               // 总伤害输出
	TotalDamageDoneToChampions    int `json:"totalDamageDoneToChampions"`    // 对英雄的总伤害输出
	TotalDamageTaken              int `json:"totalDamageTaken"`              // 承受的总伤害
	TrueDamageDone                int `json:"trueDamageDone"`                // 真实伤害输出
	TrueDamageDoneToChampions     int `json:"trueDamageDoneToChampions"`     // 对英雄的真实伤害输出
	TrueDamageTaken               int `json:"trueDamageTaken"`               // 承受的真实伤害
}

type Event struct {
	Level                   int             `json:"level"`                             // 等级
	LevelUpType             string          `json:"levelUpType,omitempty"`             // 升级类型
	ParticipantId           int             `json:"participantId"`                     // 参与者 ID
	SkillSlot               int             `json:"skillSlot,omitempty"`               // 技能选择
	ItemId                  int             `json:"itemId,omitempty"`                  // 物品 ID
	AfterId                 int             `json:"afterId,omitempty"`                 // 撤销后的物品 ID
	BeforeId                int             `json:"beforeId,omitempty"`                // 撤销前的物品 ID
	GoldGain                int             `json:"goldGain,omitempty"`                // 获得的金币
	AssistingParticipantIds []int           `json:"assistingParticipantIds,omitempty"` // 协助参与者的 ID
	Bounty                  int             `json:"bounty,omitempty"`                  // 赏金
	KillStreakLength        int             `json:"killStreakLength,omitempty"`        // 连续击杀次数
	KillerId                int             `json:"killerId,omitempty"`                // 击杀者 ID
	VictimId                int             `json:"victimId,omitempty"`                // 受害者 ID
	KillType                string          `json:"killType"`                          // 击杀类型
	Position                *Position       `json:"position,omitempty"`                // 位置
	RealTimestamp           int64           `json:"realTimestamp"`                     // 现实时间戳
	ShutdownBounty          int             `json:"shutdownBounty,omitempty"`          // 关闭赏金
	GameTimestamp           int             `json:"timestamp"`                         // 游戏时间戳
	Type                    string          `json:"type"`                              // 事件类型
	VictimDamageDealt       []*VictimDamage `json:"victimDamageDealt,omitempty"`       // 受害者的伤害信息
	VictimDamageReceived    []*VictimDamage `json:"victimDamageReceived,omitempty"`    // 受害者的承受伤害信息
	
}

type VictimDamage struct {
	Basic          bool   `json:"basic"`          // 是否基础攻击
	MagicDamage    int    `json:"magicDamage"`    // 魔法伤害
	Name           string `json:"name"`           // 名称
	ParticipantId  int    `json:"participantId"`  // 参与者 ID
	PhysicalDamage int    `json:"physicalDamage"` // 物理伤害
	SpellName      string `json:"spellName"`      // 法术名称
	SpellSlot      int    `json:"spellSlot"`      // 法术槽位
	TrueDamage     int    `json:"trueDamage"`     // 真实伤害
	Type           string `json:"type"`           // 类型
}

type Position struct {
	X int `json:"x"` // X 坐标
	Y int `json:"y"` // Y 坐标
}
