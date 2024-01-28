package test

import (
	"testing"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/service/updater"
)

func Test_update_champion(t *testing.T) {
	u := updater.NewRiotUpdater(
		updater.WithLifeTime(time.Hour * 24 * 30 * 2), // 2 month
	)

	u.UpdateAll("13.3.1")
}
