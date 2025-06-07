package util

import (
	"errors"
	"log"
	"time"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	"github.com/pixperk/notifly/notification/dispatcher"
)

func StartWorkerPool(queue <-chan common.NotificationEvent, retryQueue chan<- common.NotificationEvent, n int, cfg notification.Config) {
	for i := range n {
		go func(workerID int) {
			for job := range queue {
				// Process the job
				processJob(job, retryQueue, cfg)
			}
		}(i)
	}
}

func StartRetryWorker(retryQueue <-chan common.NotificationEvent, mainQueue chan<- common.NotificationEvent) {
	for job := range retryQueue {
		delay := time.Duration(2 * time.Second)
		time.Sleep(delay)

		log.Printf("Retrying job: %v after %v", job.NotificationId, delay)
		mainQueue <- job
	}
}

func processJob(job common.NotificationEvent, retryQueue chan<- common.NotificationEvent, cfg notification.Config) error {
	log.Printf("Processing job: %v", job.NotificationId)

	notifDispatcher := dispatcher.GetDispatcher(job, cfg)
	if notifDispatcher == nil {
		return errors.New("no dispatcher found for job type: " + job.Type)
	}

	err := notifDispatcher.Send(job)
	if err != nil {
		log.Printf("Failed to send notification for job %v: %v", job.NotificationId, err)
		job.RetryCount++
		maxRetries := getMaxRetriesByType(job)
		if job.RetryCount <= maxRetries {
			log.Printf("Retrying job %v, attempt %d/%d", job.NotificationId, job.RetryCount, maxRetries)
			retryQueue <- job
		} else {
			log.Printf("Max retries reached for job %v, giving up", job.NotificationId)
			job.RetryCount = 0 // Reset retry count for next processing
		}

		return err
	}

	log.Printf("Successfully sent notification for job %v", job.NotificationId)
	return nil
}

func getMaxRetriesByType(notificationEvent common.NotificationEvent) int {
	switch notificationEvent.Type {
	case "EMAIL":
		return 3
	case "SMS":
		return 5
	default:
		return 3
	}
}
