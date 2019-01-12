package image_url

import "context"

type urlKey struct{}

func NewContext(ctx context.Context, url URL) context.Context {
	return context.WithValue(ctx, urlKey{}, url)
}

func FromContext(ctx context.Context) (url URL, ok bool) {
	url, ok = ctx.Value(urlKey{}).(URL)
	return
}
