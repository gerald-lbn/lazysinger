package tasks

import "github.com/gerald-lbn/refrain/pkg/services"

// Register registers all task queues with the task client.
func Register(c *services.Container) {
	c.Tasks.Register(NewDownloadLyricsTaskQueue(c))
}
