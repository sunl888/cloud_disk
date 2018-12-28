package db_store

import (
	"context"
	"github.com/jinzhu/gorm"
)

type dbKey struct{}

func NewDBContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

func FromDBContext(ctx context.Context) *gorm.DB {
	return ctx.Value(dbKey{}).(*gorm.DB)
}
