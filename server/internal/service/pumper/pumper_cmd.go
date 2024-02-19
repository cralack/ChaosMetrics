package pumper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var id, loc, end, que string
var ques []riotmodel.QUECODE
var tier riotmodel.TIER
var rank uint

var Cmd = &cobra.Command{
	Use:   "pump",
	Short: "start a local pumper all",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ConvertEndmark(end); err != nil {
			global.ChaLogger.Error("init endmark failed", zap.Error(err))
			return
		}
		switch que {
		case "solo":
			ques = append(ques, riotmodel.RANKED_SOLO_5x5)
		case "flex":
			ques = append(ques, riotmodel.RANKED_FLEX_SR)
		case "all":
			ques = append(ques, riotmodel.RANKED_SOLO_5x5, riotmodel.RANKED_FLEX_SR)
		default:
			global.ChaLogger.Error("init ques failed")
			return
		}
		pumper, err := NewPumper(
			id,
			WithLoc(utils.ConvertLocStrToLocation(loc)),
			WithEndMark(tier, rank),
			WithQues(ques...),
		)
		if err != nil {
			global.ChaLogger.Error("init pumper failed", zap.Error(err))
			return
		}
		pumper.StartEngine()
		pumper.UpdateAll()
	},
}

func init() {
	Cmd.Flags().StringVar(&id, "id", "command", "set id")
	Cmd.Flags().StringVar(&loc, "loc", "na1", "set location")
	Cmd.Flags().StringVar(&end, "end", "d1", "set end mark")
	Cmd.Flags().StringVar(&que, "que", "all", "set que")
}

func ConvertEndmark(str string) (err error) {
	str = strings.ToUpper(str)
	rankStr := str[len(str)-1]
	rank = uint(rankStr - '0')
	if 4 < rank || rank < 0 {
		return errors.New("wrong rank")
	}

	switch str[0] {
	case 'C':
		tier = riotmodel.CHALLENGER
	case 'G':
		tier = riotmodel.GRANDMASTER
	case 'M':
		tier = riotmodel.MASTER
	case 'D':
		tier = riotmodel.DIAMOND
	case 'E':
		tier = riotmodel.EMERALD
	case 'P':
		tier = riotmodel.PLATINUM
	// case 'G':
	// 	tier = riotmodel.GOLD
	case 'S':
		tier = riotmodel.SILVER
	case 'B':
		tier = riotmodel.BRONZE
	case 'I':
		tier = riotmodel.IRON
	default:
		fmt.Println("Unknown tier")
	}
	return nil
}
