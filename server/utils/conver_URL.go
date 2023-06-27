package utils

// LOCATION CODE
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

func ConvertPlatformURL(loc uint) string {
	platformHosts := map[uint]string{
		BR1:  "br1.api.riotgames.com",
		EUN1: "eun1.api.riotgames.com",
		EUW1: "euw1.api.riotgames.com",
		JP1:  "jp1.api.riotgames.com",
		KR1:  "kr.api.riotgames.com",
		LA1:  "la1.api.riotgames.com",
		LA2:  "la2.api.riotgames.com",
		NA1:  "na1.api.riotgames.com",
		OC1:  "oc1.api.riotgames.com",
		PH2:  "ph2.api.riotgames.com",
		RU:   "ru.api.riotgames.com",
		SG2:  "sg2.api.riotgames.com",
		TH2:  "th2.api.riotgames.com",
		TR1:  "tr1.api.riotgames.com",
		TW2:  "tw2.api.riotgames.com",
		VN2:  "vn2.api.riotgames.com",
	}
	return platformHosts[loc]
}

func ConvertPlatformToHost(loc uint) string {
	platformToRegion := map[uint]string{
		// The AMERICAS routing value serves NA, BR, LAN and LAS.
		BR1: "AMERICAS",
		LA1: "AMERICAS",
		LA2: "AMERICAS",
		NA1: "AMERICAS",
		// The ASIA routing value serves KR and JP.
		KR1: "ASIA",
		JP1: "ASIA",
		// The EUROPE routing value serves EUNE, EUW, TR and RU.
		EUN1: "EUROPE",
		EUW1: "EUROPE",
		TR1:  "EUROPE",
		RU:   "EUROPE",
		// The SEA routing value serves OCE, PH2, SG2, TH2, TW2 and VN2.
		OC1: "SEA",
		PH2: "SEA",
		SG2: "SEA",
		TH2: "SEA",
		TW2: "SEA",
		VN2: "SEA",
	}

	regionToHost := map[string]string{
		"AMERICAS": "americas.api.riotgames.com",
		"ASIA":     "asia.api.riotgames.com",
		"EUROPE":   "europe.api.riotgames.com",
		"SEA":      "sea.api.riotgames.com",
	}

	region := platformToRegion[loc]
	host := regionToHost[region]
	return host
}
