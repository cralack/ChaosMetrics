package test

import (
	"testing"
	"time"

	"github.com/cralack/ChaosMetrics/server/service/updater"
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
	u.UpdatePerks()
	// u.UpdateItems(versions[18])
	// u.UpdateChampions(versions[18])
	for _, ver := range versions {
		if ver[:2] == "13" {
			logger.Info(ver)
			u.UpdateItems(ver)
			u.UpdateChampions(ver)
		} else {
			break
		}
	}
}
