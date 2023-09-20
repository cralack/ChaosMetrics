package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/service/updater"
)

func Test_update_champion(t *testing.T) {
	var (
		versions []string
	)
	u := updater.NewRiotUpdater()

	if u.CurVersion == "" {
		versions = u.UpdateVersions()
		u.CurVersion = versions[0]
	}
	u.UpdatePerks()
	// u.UpdateItems(u.CurVersion)
	// u.UpdateChampions(u.CurVersion)
	for _, ver := range versions {
		if ver[:2] == "13" {
			u.UpdateItems(ver)
			u.UpdateChampions(ver)
		} else {
			break
		}
	}
}
