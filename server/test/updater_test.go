package test

import (
	"testing"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/service/updater"
)

func Test_update_champion(t *testing.T) {
	var (
		versions []string
	)
	u := updater.NewRiotUpdater(
		updater.WithLifeTime(time.Hour * 24 * 30 * 2), // 2 month
	)

	if u.CurVersion == "" {
		versions = u.UpdateVersions()
	}
	t.Log(len(versions))
	for _, ver := range versions {
		if ver[:2] >= "13" {
			logger.Info(ver)
			u.UpdatePerks(ver)
			u.UpdateItems(ver)
			u.UpdateChampions(ver)
		} else {
			break
		}
	}
}
