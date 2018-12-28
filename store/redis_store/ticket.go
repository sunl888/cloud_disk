package redis_store

import (
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"github.com/wq1019/cloud_disk/model"
)

type redisTicket struct {
	client *redis.Client
}

func (rt *redisTicket) id2key(id string) string {
	return "ticket:" + id
}

func (rt *redisTicket) TicketIsNotExistErr(err error) bool {
	return model.TicketIsNotExistErr(err)
}

func (rt *redisTicket) TicketLoad(id string) (ticket *model.Ticket, err error) {

	if id == "" {
		return nil, model.ErrTicketNotExist
	}

	res, err := rt.client.Get(rt.id2key(id)).Result()
	if err != nil {
		if err == redis.Nil {
			err = model.ErrTicketNotExist
		}
		return nil, err
	}
	ticket = &model.Ticket{}
	if err = msgpack.Unmarshal([]byte(res), ticket); err != nil {
		return nil, err
	}
	return
}

func (rt *redisTicket) TicketCreate(ticket *model.Ticket) error {
	key := rt.id2key(ticket.Id)
	if res, err := rt.client.Exists(key).Result(); err != nil {
		return err
	} else if res != 0 {
		return model.ErrTicketExisted
	}

	b, err := msgpack.Marshal(ticket)
	if err != nil {
		return err
	}
	return rt.client.Set(key, b, ticket.ExpiredAt.Sub(ticket.CreatedAt)).Err()
}

func (rt *redisTicket) TicketDelete(id string) error {
	return rt.client.Del(rt.id2key(id)).Err()
}

func NewRedisTicket(client *redis.Client) model.TicketStore {
	return &redisTicket{client: client}
}
