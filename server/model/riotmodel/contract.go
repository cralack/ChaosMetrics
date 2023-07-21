package riotmodel

// DEFINE BY PACK JSON

// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
// type Unmarshaler interface {
// 	UnmarshalJSON([]byte) error
// }

type DTO interface {
	UnmarshalJSON(data []byte) error
}

// TIER LEVEL
const (
	CHALLENGER  = iota // top 200 players
	GRANDMASTER        // top 201-701 players
	MASTER             // top 702-2922 players
	DIAMOND
	PLATINUM
	GOLD
	SILVER
	BRONZE
	IRON
)

const (
	LOC_ALL = iota
	LOC_AMERICAS
	LOC_ASIA
	LOC_EUROPE
	LOC_SEA
)

// 16 LOCATION CODE
const (
	BR1  = iota // 巴西
	EUN1        // 欧洲东北
	EUW1        // 欧洲西部
	JP1         // 日本
	KR1         // 韩国
	LA1         // 拉丁美洲北部
	LA2         // 拉丁美洲南部
	NA1         // 北美洲
	OC1         // 大洋洲
	PH2         // 菲律宾
	RU          // 俄罗斯
	SG2         // 新加坡
	TH2         // 泰国
	TR1         // 土耳其
	TW2         // 台湾
	VN2         // 越南
)

// QUE CODE
const (
	RANKED_SOLO_5x5 = iota
	RANKED_FLEX_SR
	RANKED_FLEX_TT
)

// 定义常量
const (
	LANG_cs_CZ = iota // Czech (Czech Republic)
	LANG_el_GR        // Greek (Greece)
	LANG_pl_PL        // Polish (Poland)
	LANG_ro_RO        // Romanian (Romania)
	LANG_hu_HU        // Hungarian (Hungary)
	LANG_en_GB        // English (United Kingdom)
	LANG_de_DE        // German (Germany)
	LANG_es_ES        // Spanish (Spain)
	LANG_it_IT        // Italian (Italy)
	LANG_fr_FR        // French (France)
	LANG_ja_JP        // Japanese (Japan)
	LANG_ko_KR        // Korean (Korea)
	LANG_es_MX        // Spanish (Mexico)
	LANG_es_AR        // Spanish (Argentina)
	LANG_pt_BR        // Portuguese (Brazil)
	LANG_en_US        // English (United States)
	LANG_en_AU        // English (Australia)
	LANG_ru_RU        // Russian (Russia)
	LANG_tr_TR        // Turkish (Turkey)
	LANG_ms_MY        // Malay (Malaysia)
	LANG_en_PH        // English (Republic of the Philippines)
	LANG_en_SG        // English (Singapore)
	LANG_th_TH        // Thai (Thailand)
	LANG_vi_VN        // Vietnamese (Viet Nam)
	LANG_id_ID        // Indonesian (Indonesia)
	LANG_zh_MY        // Chinese (Malaysia)
	LANG_zh_CN        // Chinese (China)
	LANG_zh_TW        // Chinese (Taiwan)
)
