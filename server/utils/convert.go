package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func ConvertLocationToLoHoSTR(loc riotmodel.LOCATION) (loCode, hostURL string) {
	platformCodeMap := map[riotmodel.LOCATION]string{
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

func ConvertLocStrToLocation(loCode string) riotmodel.LOCATION {
	platformCodeMap := map[string]riotmodel.LOCATION{
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

func ConvertLocationToRegionHost(loc riotmodel.LOCATION) string {
	platformToRegion := map[riotmodel.LOCATION]string{
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

func ConvertRegionStrToArea(region string) riotmodel.AREA {
	regionCode := map[string]riotmodel.AREA{
		"AMERICAS": riotmodel.LOC_AMERICAS,
		"ASIA":     riotmodel.LOC_ASIA,
		"EUROPE":   riotmodel.LOC_EUROPE,
		"SEA":      riotmodel.LOC_SEA,
	}
	if region == "" {
		return regionCode["AMERICA"]
	}
	return regionCode[region]
}

func ConvertLangToLangStr(lang riotmodel.LANG) string {
	langMap := map[riotmodel.LANG]string{
		riotmodel.LANG_cs_CZ: "cs_CZ",
		riotmodel.LANG_el_GR: "el_GR",
		riotmodel.LANG_pl_PL: "pl_PL",
		riotmodel.LANG_ro_RO: "ro_RO",
		riotmodel.LANG_hu_HU: "hu_HU",
		riotmodel.LANG_en_GB: "en_GB",
		riotmodel.LANG_de_DE: "de_DE",
		riotmodel.LANG_es_ES: "es_ES",
		riotmodel.LANG_it_IT: "it_IT",
		riotmodel.LANG_fr_FR: "fr_FR",
		riotmodel.LANG_ja_JP: "ja_JP",
		riotmodel.LANG_ko_KR: "ko_KR",
		riotmodel.LANG_es_MX: "es_MX",
		riotmodel.LANG_es_AR: "es_AR",
		riotmodel.LANG_pt_BR: "pt_BR",
		riotmodel.LANG_en_US: "en_US",
		riotmodel.LANG_en_AU: "en_AU",
		riotmodel.LANG_ru_RU: "ru_RU",
		riotmodel.LANG_tr_TR: "tr_TR",
		riotmodel.LANG_ms_MY: "ms_MY",
		riotmodel.LANG_en_PH: "en_PH",
		riotmodel.LANG_en_SG: "en_SG",
		riotmodel.LANG_th_TH: "th_TH",
		riotmodel.LANG_vi_VN: "vi_VN",
		riotmodel.LANG_id_ID: "id_ID",
		riotmodel.LANG_zh_MY: "zh_MY",
		riotmodel.LANG_zh_CN: "zh_CN",
		riotmodel.LANG_zh_TW: "zh_TW",
	}

	return langMap[lang]
}

func ConvertSliceToStr(s []string) string {
	return strings.Join(s, ",")
}

func ConvertStrToSlice(str string) []string {
	return strings.Split(str, ",")
}

// ConvertVersionToIdx return a 4 digit str xxyy (1401)
func ConvertVersionToIdx(version string) (uint, error) {
	versions := strings.Split(version, ".")
	if len(versions) < 3 {
		return 0, errors.New("wrong version")
	}
	versionNums := make([]int, 0, 2)
	for _, str := range versions[:2] {
		if num, err := strconv.Atoi(str); err != nil {
			return 0, err
		} else {
			versionNums = append(versionNums, num)
		}
	}

	str := fmt.Sprintf("%02d%02d", versionNums[0], versionNums[1])
	idx, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return uint(idx), nil
}

func ConvertQueToQueSTR(que riotmodel.QUECODE) string {
	switch que {
	case riotmodel.RANKED_SOLO_5x5:
		return "RANKED_SOLO_5x5"
	case riotmodel.RANKED_FLEX_SR:
		return "RANKED_FLEX_SR"
	default:
		return ""
	}
}

// ConvertQueStrToQue
func _(que string) riotmodel.QUECODE {
	switch que {
	case "RANKED_SOLO_5x5":
		return riotmodel.RANKED_SOLO_5x5
	case "RANKED_FLEX_SR":
		return riotmodel.RANKED_FLEX_SR
	default:
		return 999
	}
}

func GetCurMajorVersions() []string {
	res := global.ChaRDB.HGet(context.Background(), "/version", "versions")
	if res.Err() != nil {
		global.ChaLogger.Error(res.Err().Error())
		return []string{}
	}
	versions := make([]string, 0)
	if err := json.Unmarshal([]byte(res.Val()), &versions); err != nil {
		global.ChaLogger.Error(err.Error())
		return []string{}
	}
	majorVersion, _ := strconv.Atoi(versions[0][:2])
	majorVersion *= 100
	for i, v := range versions {
		if ver, _ := ConvertVersionToIdx(v); int(ver) <= majorVersion {
			versions = versions[:i]
			break
		}
	}
	return versions
}
