package riotmodel

type TIER uint

const (
	CHALLENGER  TIER = iota // most top 200 players
	GRANDMASTER             // most top 201-701 players
	MASTER                  // most top 702-2922 players
	DIAMOND
	EMERALD
	PLATINUM
	GOLD
	SILVER
	BRONZE
	IRON
)

type AREA uint

const (
	LOC_ALL AREA = iota
	LOC_AMERICAS
	LOC_ASIA
	LOC_EUROPE
	LOC_SEA
)

type LOCATION uint

const (
	BR1  LOCATION = iota + 1 // 巴西
	EUN1                     // 欧洲东北
	EUW1                     // 欧洲西部
	JP1                      // 日本
	KR1                      // 韩国
	LA1                      // 拉丁美洲北部
	LA2                      // 拉丁美洲南部
	NA1                      // 北美洲
	OC1                      // 大洋洲
	PH2                      // 菲律宾
	RU                       // 俄罗斯
	SG2                      // 新加坡
	TH2                      // 泰国
	TR1                      // 土耳其
	TW2                      // 台湾
	VN2                      // 越南
)

type QUECODE uint

const (
	RANKED_SOLO_5x5 QUECODE = iota // 单排/双排
	RANKED_FLEX_SR                 // 灵活排位
	// RANKED_FLEX_TT // abandoned?
)

type LANG uint

const (
	LANG_cs_CZ LANG = iota // Czech (Czech Republic)
	LANG_el_GR             // Greek (Greece)
	LANG_pl_PL             // Polish (Poland)
	LANG_ro_RO             // Romanian (Romania)
	LANG_hu_HU             // Hungarian (Hungary)
	LANG_en_GB             // English (United Kingdom)
	LANG_de_DE             // German (Germany)
	LANG_es_ES             // Spanish (Spain)
	LANG_it_IT             // Italian (Italy)
	LANG_fr_FR             // French (France)
	LANG_ja_JP             // Japanese (Japan)
	LANG_ko_KR             // Korean (Korea)
	LANG_es_MX             // Spanish (Mexico)
	LANG_es_AR             // Spanish (Argentina)
	LANG_pt_BR             // Portuguese (Brazil)
	LANG_en_US             // English (United States)
	LANG_en_AU             // English (Australia)
	LANG_ru_RU             // Russian (Russia)
	LANG_tr_TR             // Turkish (Turkey)
	LANG_ms_MY             // Malay (Malaysia)
	LANG_en_PH             // English (Republic of the Philippines)
	LANG_en_SG             // English (Singapore)
	LANG_th_TH             // Thai (Thailand)
	LANG_vi_VN             // Vietnamese (Viet Nam)
	LANG_id_ID             // Indonesian (Indonesia)
	LANG_zh_MY             // Chinese (Malaysia)
	LANG_zh_CN             // Chinese (China)
	LANG_zh_TW             // Chinese (Taiwan)
)

type GAMEMODE uint

const (
	CLASSIC GAMEMODE = iota // 召唤师峡谷 mapid:11
	ARAM                    // 大乱斗 mapid:12
	CHERRY                  // 斗魂竞技场 mapid:30
	// ONEFORALL                 // 镜像模式 mapid:11
	// NEXUSBLITZ                  // 扭曲丛林？ mapid:21
	// CONVERGENCE                 // ??? mapid:22
)
