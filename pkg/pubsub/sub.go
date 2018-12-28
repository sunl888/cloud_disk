package pubsub

import (
	"context"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type SubQueue interface {
	Channel() string
	Process(ctx context.Context, message string)
}

type Sub struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
	queueNum    int
	subs        map[string][]SubQueue
	execChan    chan *redis.Message
}

func (bq *Sub) Sub(ctx context.Context) {
	channels := make([]string, 0, len(bq.subs))
	for channel := range bq.subs {
		channels = append(channels, channel)
	}
	pubsub := bq.RedisClient.Subscribe(channels...)
	defer func() {
		close(bq.execChan)
		pubsub.Close()
	}()
	//defer func() {
	//	recover() // fix #2480
	//}()
	for i := 0; i < bq.queueNum; i++ {
		go bq.process(ctx)
	}
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			bq.Logger.Error("receive message error.", zap.Error(err))
		}
		bq.execChan <- msg
	}
}

func (bq *Sub) process(ctx context.Context) {
	for {
		select {
		case msg, ok := <-bq.execChan:
			if !ok {
				return
			}
			subs, ok := bq.subs[msg.Channel]
			if ok {
				for _, sub := range subs {
					sub.Process(ctx, msg.Payload)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (bq *Sub) RegisterSub(sqs ...SubQueue) {
	for _, sq := range sqs {
		channel := sq.Channel()
		_, ok := bq.subs[channel]
		if ok {
			bq.subs[channel] = append(bq.subs[channel], sq)
		} else {
			bq.subs[channel] = []SubQueue{sq}
		}
	}
}

func NewSub(redisClient *redis.Client, logger *zap.Logger, queueNum int) *Sub {
	execChan := make(chan *redis.Message, queueNum)
	subs := make(map[string][]SubQueue)
	return &Sub{RedisClient: redisClient, Logger: logger, subs: subs, execChan: execChan, queueNum: queueNum}
}
