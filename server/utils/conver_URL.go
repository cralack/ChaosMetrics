package utils

import (
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func ConvertPlatformURL(loc uint) string {
	platformHosts := map[uint]string{
		riotmodel.BR1:  "br1.api.riotgames.com",
		riotmodel.EUN1: "eun1.api.riotgames.com",
		riotmodel.EUW1: "euw1.api.riotgames.com",
		riotmodel.JP1:  "jp1.api.riotgames.com",
		riotmodel.KR1:  "kr.api.riotgames.com",
		riotmodel.LA1:  "la1.api.riotgames.com",
		riotmodel.LA2:  "la2.api.riotgames.com",
		riotmodel.NA1:  "na1.api.riotgames.com",
		riotmodel.OC1:  "oc1.api.riotgames.com",
		riotmodel.PH2:  "ph2.api.riotgames.com",
		riotmodel.RU:   "ru.api.riotgames.com",
		riotmodel.SG2:  "sg2.api.riotgames.com",
		riotmodel.TH2:  "th2.api.riotgames.com",
		riotmodel.TR1:  "tr1.api.riotgames.com",
		riotmodel.TW2:  "tw2.api.riotgames.com",
		riotmodel.VN2:  "vn2.api.riotgames.com",
	}
	return "https://" + platformHosts[loc]
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
