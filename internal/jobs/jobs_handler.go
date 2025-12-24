package jobs

import (
	"github.com/techies/streamify/internal/app"
)

func StartAllJobs(appCfg *app.AppConfig) {
	StartUserCleanupJob(appCfg)
}
