package utils

import (
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func ConvertHostURL(loc uint) (loCode, hostURL string) {
	platformCodeMap := map[uint]string{
		riotmodel.BR1:  "br1",
		riotmodel.EUN1: "eun1",
		riotmodel.EUW1: "euw1",
		riotmodel.JP1:  "jp1",
		riotmodel.KR1:  "kr",
		riotmodel.LA1:  "la1",
		riotmodel.LA2:  "la2",
		riotmodel.NA1:  "na1",
		riotmodel.OC1:  "oc1",
		riotmodel.PH2:  "ph2",
		riotmodel.RU:   "ru",
		riotmodel.SG2:  "sg2",
		riotmodel.TH2:  "th2",
		riotmodel.TR1:  "tr1",
		riotmodel.TW2:  "tw2",
		riotmodel.VN2:  "vn2",
	}
	loCode = platformCodeMap[loc]
	host := ".api.riotgames.com"
	return loCode, "https://" + loCode + host
}
func ConverHostLoCode(loCode string) uint {
	platformCodeMap := map[string]uint{
		"br1":  riotmodel.BR1,
		"eun1": riotmodel.EUN1,
		"euw1": riotmodel.EUW1,
		"jp1":  riotmodel.JP1,
		"kr":   riotmodel.KR1,
		"la1":  riotmodel.LA1,
		"la2":  riotmodel.LA2,
		"na1":  riotmodel.NA1,
		"oc1":  riotmodel.OC1,
		"ph2":  riotmodel.PH2,
		"ru":   riotmodel.RU,
		"sg2":  riotmodel.SG2,
		"th2":  riotmodel.TH2,
		"tr1":  riotmodel.TR1,
		"tw2":  riotmodel.TW2,
		"vn2":  riotmodel.VN2,
	}
	
	return platformCodeMap[loCode]
}

func ConvertPlatformToHost(loc uint) string {
	platformToRegion := map[uint]string{
		// The AMERICAS routing value serves NA, BR, LAN and LAS.
		riotmodel.BR1: "AMERICAS",
		riotmodel.LA1: "AMERICAS",
		riotmodel.LA2: "AMERICAS",
		riotmodel.NA1: "AMERICAS",
		// The ASIA routing value serves KR and JP.
		riotmodel.KR1: "ASIA",
		riotmodel.JP1: "ASIA",
		// The EUROPE routing value serves EUNE, EUW, TR and RU.
		riotmodel.EUN1: "EUROPE",
		riotmodel.EUW1: "EUROPE",
		riotmodel.TR1:  "EUROPE",
		riotmodel.RU:   "EUROPE",
		// The SEA routing value serves OCE, PH2, SG2, TH2, TW2 and VN2.
		riotmodel.OC1: "SEA",
		riotmodel.PH2: "SEA",
		riotmodel.SG2: "SEA",
		riotmodel.TH2: "SEA",
		riotmodel.TW2: "SEA",
		riotmodel.VN2: "SEA",
	}

	regionToHost := map[string]string{
		"AMERICAS": "americas.api.riotgames.com",
		"ASIA":     "asia.api.riotgames.com",
		"EUROPE":   "europe.api.riotgames.com",
		"SEA":      "sea.api.riotgames.com",
	}

	region := platformToRegion[loc]
	host := regionToHost[region]
	return "https://" + host
}
