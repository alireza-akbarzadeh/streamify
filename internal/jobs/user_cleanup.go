package jobs

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/techies/streamify/internal/app"
)

// StartUserCleanupJob schedules a daily job to permanently delete old soft-deleted users
func StartUserCleanupJob(app *app.AppConfig) {
	c := cron.New()
	_, err := c.AddFunc("0 3 * * *", func() {
		ctx := context.Background()
		if err := app.DB.PermanentlyDeleteOldSoftDeletedUsers(ctx); err != nil {
			log.Printf("User cleanup job failed: %v", err)
		} else {
			log.Println("User cleanup job completed successfully")
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule user cleanup job: %v", err)
	}

	c.Start()
}
