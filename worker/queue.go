package worker

import (
	"github.com/terorie/yt-mango/store"
	log "github.com/sirupsen/logrus"
	"time"
)

// Queue handler:
// Reads and writes to queue in the background

// TODO Handle errors
func (c *workerContext) handleQueueWrite() {
	var timeOut <-chan time.Time
	for { select {
		case <-c.ctxt.Done():
			timeOut = time.After(1 * time.Second)

		case id := <-c.resultIDs:
			err := store.DoneVideoID(id)
			if err != nil {
				log.Errorf("Marking video \"%s\" as done failed: %s", id, err.Error())
				c.errors <- err
			}

		case id := <-c.failIDs:
			err := store.FailedVideoID(id)
			if err != nil {
				log.Errorf("Marking video \"%s\" as failed failed: %s", id, err.Error())
				c.errors <- err
			}

		case id := <-c.newIDs:
			err := store.SubmitVideoID(id)
			if err != nil {
				log.Errorf("Pushing related video IDs of video \"%s\" failed: %s", id, err.Error())
				c.errors <- err
			}

		case <-timeOut:
			return
	}}
}

func (c *workerContext) handleQueueReceive() {
	for { select {
		case <-c.ctxt.Done():
			return

		default:
			videoId, err := store.GetScheduledVideoID()
			if err != nil && err.Error() != "redis: nil" {
				log.Error("Queue error: ", err.Error())
				c.errors <- err
			}
			if videoId == "" {
				// Queue is empty
				log.Info("No jobs on queue. Waiting 1 second.")
				time.Sleep(1 * time.Second)
				continue
			}

			c.jobs <- videoId
	}}
}