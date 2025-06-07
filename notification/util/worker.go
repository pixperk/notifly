package util

import (
	"errors"
	"log"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	"github.com/pixperk/notifly/notification/dispatcher"
)

func StartWorkerPool(queue <-chan common.NotificationEvent, n int, cfg notification.Config) {
	for i := range n {
		go func(workerID int) {
			for job := range queue {
				// Process the job
				processJob(job, cfg)
			}
		}(i)
	}
}

func processJob(job common.NotificationEvent, cfg notification.Config) error {
	log.Printf("Processing job: %v", job.NotificationId)

	notifDispatcher := dispatcher.GetDispatcher(job, cfg)
	if notifDispatcher == nil {
		return errors.New("no dispatcher found for job type: " + job.Type)
	}

	err := notifDispatcher.Send(job)
	if err != nil {
		log.Printf("Failed to send notification for job %v: %v", job.NotificationId, err)
		//TODO: Handle retry logic or error handling as needed
		return err
	}

	log.Printf("Successfully sent notification for job %v", job.NotificationId)
	return nil
}
