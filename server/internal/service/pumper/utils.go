package pumper

import (
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func getQueueString(que riotmodel.QUECODE) string {
	switch que {
	case riotmodel.RANKED_SOLO_5x5:
		return "RANKED_SOLO_5x5"
	case riotmodel.RANKED_FLEX_SR:
		return "RANKED_FLEX_SR"
	case riotmodel.RANKED_FLEX_TT:
		return "RANKED_FLEX_TT"
	default:
		return ""
	}
}

func ConvertRankToStr(tier riotmodel.TIER, div uint) (string, string) {
	var d string
	switch div {
	case 1:
		d = "I"
	case 2:
		d = "II"
	case 3:
		d = "III"
	case 4:
		d = "IV"
	}
	switch tier {
	case riotmodel.CHALLENGER:
		return "CHALLENGER", "I"
	case riotmodel.GRANDMASTER:
		return "GRANDMASTER", "I"
	case riotmodel.MASTER:
		return "MASTER", "I"
	case riotmodel.DIAMOND:
		return "DIAMOND", d
	case riotmodel.EMERALD:
		return "EMERALD", d
	case riotmodel.PLATINUM:
		return "PLATINUM", d
	case riotmodel.GOLD:
		return "GOLD", d
	case riotmodel.SILVER:
		return "SILVER", d
	case riotmodel.BRONZE:
		return "BRONZE", d
	case riotmodel.IRON:
		return "IRON", d
	}

	return "", ""
}

func ConvertStrToRank(tierStr, divStr string) (riotmodel.TIER, uint) {
	var tier riotmodel.TIER
	var div uint

	switch tierStr {
	case "CHALLENGER":
		tier = riotmodel.CHALLENGER
	case "GRANDMASTER":
		tier = riotmodel.GRANDMASTER
	case "MASTER":
		tier = riotmodel.MASTER
	case "DIAMOND":
		tier = riotmodel.DIAMOND
	case "EMERALD":
		tier = riotmodel.EMERALD
	case "PLATINUM":
		tier = riotmodel.PLATINUM
	case "GOLD":
		tier = riotmodel.GOLD
	case "SILVER":
		tier = riotmodel.SILVER
	case "BRONZE":
		tier = riotmodel.BRONZE
	case "IRON":
		tier = riotmodel.IRON
	default:
		return 0, 0
	}

	switch divStr {
	case "I":
		div = 1
	case "II":
		div = 2
	case "III":
		div = 3
	case "IV":
		div = 4
	default:
		return 0, 0
	}

	return tier, div
}
