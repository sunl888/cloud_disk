package pubsub

import "context"

type pubKey struct{}

func NewContext(ctx context.Context, p PubQueue) context.Context {
	return context.WithValue(ctx, pubKey{}, p)
}

func FromContext(ctx context.Context) PubQueue {
	return ctx.Value(pubKey{}).(PubQueue)
}
