package service

import (
	"log"

	"github.com/robfig/cron/v3"
)

func SyncScopus(personnelService PersonnelService) {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 0 12 * * 0", func() {
		log.Println("[CRON] Start syncing research from Scopus")

		count, err := personnelService.SyncAllFromScopus()
		if err != nil {
			log.Println("[CRON] Sync failed:", err)
			return
		}

		log.Printf("[CRON] Sync completed, processed %d personnels\n", count)
	})

	if err != nil {
		log.Fatal("cannot start scopus cron:", err)
	}

	c.Start()
}
