package pumper

import (
	"fmt"

	"github.com/cralack/ChaosMetrics/server/internal/service/master"
)

func (p *Pumper) TaskHandlers(body []byte) error {
	task, err := master.Decode(body)
	if err != nil {
		return err
	}
	// todo：handler
	// loc := utils.ConvertLocStrToLocation(task.Loc)
	// switch task.Type {
	// case EntryTypeKey:
	// 	err = p.FetchEntryByName(task.SumName, loc)
	// case MatchTypeKey:
	// 	err = p.FetchMatchByName(task.SumName, loc)
	// case SummonerTypeKey:
	// }
	p.logger.Debug(fmt.Sprintf("task info: %v,%v", task.Loc, task.ID))
	return nil
}
