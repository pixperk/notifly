package util

import (
	"log"

	"github.com/pixperk/notifly/common"
)

func StartWorkerPool(queue <-chan common.NotificationEvent, n int) {
	for i := range n {
		go func(workerID int) {
			for job := range queue {
				// Process the job
				processJob(job)
			}
		}(i)
	}
}

func processJob(job common.NotificationEvent) {
	//illustrative processing logic, will be replaced with actual logic
	log.Printf("Processing job: %v", job)

}
