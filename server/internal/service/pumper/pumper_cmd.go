package pumper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var id_, loc_, end_, que_ string
var ques []riotmodel.QUECODE
var tier_ riotmodel.TIER
var rank_ uint

var Cmd = &cobra.Command{
	Use:   "pump",
	Short: "start a local pumper all",
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd.Context())
	},
}

func init() {
	Cmd.Flags().StringVar(&id_, "id", "command", "set id")
	Cmd.Flags().StringVar(&loc_, "loc", "na1", "set location")
	Cmd.Flags().StringVar(&end_, "end", "d1", "set end mark")
	Cmd.Flags().StringVar(&que_, "que", "all", "set que")
}

func Run(ctx context.Context) {
	if err := ConvertEndmark(end_); err != nil {
		global.ChaLogger.Error("init endmark failed", zap.Error(err))
		return
	}
	switch que_ {
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
		id_,
		WithLoc(utils.ConvertLocStrToLocation(loc_)),
		WithEndMark(tier_, rank_),
		WithQues(ques...),
		WithContext(ctx),
	)
	if err != nil {
		global.ChaLogger.Error("init pumper failed", zap.Error(err))
		return
	}

	pumper.StartEngine()
	// pumper.UpdateAll()
}

func ConvertEndmark(str string) (err error) {
	str = strings.ToUpper(str)
	rankStr := str[len(str)-1]
	rank_ = uint(rankStr - '0')
	if 4 < rank_ || rank_ < 0 {
		return errors.New("wrong rank")
	}

	switch str[0] {
	case 'C':
		tier_ = riotmodel.CHALLENGER
	case 'G':
		tier_ = riotmodel.GRANDMASTER
	case 'M':
		tier_ = riotmodel.MASTER
	case 'D':
		tier_ = riotmodel.DIAMOND
	case 'E':
		tier_ = riotmodel.EMERALD
	case 'P':
		tier_ = riotmodel.PLATINUM
	// case 'G':
	// 	tier = riotmodel.GOLD
	case 'S':
		tier_ = riotmodel.SILVER
	case 'B':
		tier_ = riotmodel.BRONZE
	case 'I':
		tier_ = riotmodel.IRON
	default:
		fmt.Println("Unknown tier")
	}
	return nil
}
