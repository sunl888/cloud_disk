package image_url

import (
	"context"
	"errors"
	"strings"
)

var (
	XSOptions = []Option{Width(50), Height(50), Quality(90)}
	SMOptions = []Option{Width(100), Height(100), Quality(90)}
	MDOptions = []Option{Width(200), Height(200), Quality(90)}
	LGOptions = []Option{Width(500), Quality(90)}
)

type options struct {
	width, height               uint
	widthPercent, heightPercent float32
	quality                     uint8
}

type Option func(*options)

func Width(width uint) Option {
	return func(o *options) {
		o.width = width
	}
}

func Height(height uint) Option {
	return func(o *options) {
		o.height = height
	}
}

func WidthPercent(widthPercent float32) Option {
	return func(o *options) {
		o.widthPercent = widthPercent
	}
}

func HeightPercent(heightPercent float32) Option {
	return func(o *options) {
		o.heightPercent = heightPercent
	}
}

func Quality(quality uint8) Option {
	return func(o *options) {
		o.quality = quality
	}
}

var defaultURLOptions = options{
	quality: 90,
}

type URL interface {
	Generate(hashValue string, opt ...Option) string
}

func Generate(ctx context.Context, hashValue string, opt ...Option) (string, error) {
	url, ok := FromContext(ctx)
	if !ok {
		return "", errors.New("context中不存在 URL")
	}
	return url.Generate(hashValue, opt...), nil
}

func MustGenerate(ctx context.Context, hashValue string, opt ...Option) (string) {
	if u, err := Generate(ctx, hashValue, opt...); err != nil {
		panic(err)
	} else {
		return u
	}
}

type Hash2StorageName interface {
	Convent(hash string) (storageName string, err error)
}

type Hash2StorageNameFunc func(hash string) (storageName string, err error)

func (f Hash2StorageNameFunc) Convent(hash string) (storageName string, err error) {
	return f(hash)
}

func DefaultHash2StorageNameFunc(hash string) (storageName string, err error) {
	return hash, nil
}

func TwoCharsPrefixHash2StorageNameFunc(hash string) (storageName string, err error) {
	if len(hash) <= 2 {
		return "", errors.New("hash length must greater than 2 chars")
	}
	sb := strings.Builder{}
	sb.WriteString(hash[:2])
	sb.WriteString("/")
	sb.WriteString(hash[2:])
	return sb.String(), nil
}
