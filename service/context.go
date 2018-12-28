package service

import "context"

type serviceKey struct{}

func NewContext(ctx context.Context, s Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, s)
}

func FromContext(ctx context.Context) Service {
	return ctx.Value(serviceKey{}).(Service)
}
