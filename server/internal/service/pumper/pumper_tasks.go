package pumper

import (
	"errors"
	"fmt"

	"github.com/cralack/ChaosMetrics/server/internal/service/master"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func (p *Pumper) TaskHandlers(body []byte) error {
	task, err := master.Decode(body)
	if err != nil {
		return err
	}
	taskLoc := utils.ConvertLocStrToLocation(task.Loc)
	sumn := p.LoadSingleSummoner(task.SumName, task.Loc)
	if sumn == nil {
		return errors.New("summoner not found")
	}

	switch task.Type {
	case entryBySumnID:
		p.FetchEntryBySumnID(sumn.MetaSummonerID, taskLoc)
	case matchBySumnID:
		p.FetchMatchBySumnID(sumn.MetaSummonerID, taskLoc)
	}
	p.logger.Debug(fmt.Sprintf("task info: %v,%v", task.Loc, task.ID))
	return nil
}
