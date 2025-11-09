package worker

import "github.com/hibiken/asynq"

// NewClient returns a new Client instance given a redis connection option.
func NewClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}
