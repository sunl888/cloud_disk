package pubsub

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type PubQueue interface {
	Pub(channel, message string)
}

type Pub struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
}

func (bq *Pub) Pub(channel, message string) {
	err := bq.RedisClient.Publish(channel, message).Err()
	if err != nil {
		bq.Logger.Error("join queue failed", zap.String("channel", channel), zap.Error(err))
	}
}

func NewPub(redisClient *redis.Client, logger *zap.Logger) PubQueue {
	return &Pub{RedisClient: redisClient, Logger: logger}
}
