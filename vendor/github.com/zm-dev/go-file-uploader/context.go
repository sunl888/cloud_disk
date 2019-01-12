package go_file_uploader

import "context"

type uploaderKey struct{}

func NewContext(ctx context.Context, u Uploader) context.Context {
	return context.WithValue(ctx, uploaderKey{}, u)
}

func FromContext(ctx context.Context) (u Uploader, ok bool) {
	u, ok = ctx.Value(uploaderKey{}).(Uploader)
	return
}

