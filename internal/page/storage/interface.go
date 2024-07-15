package storage

import (
	"context"
	"projects/DAB/internal/page"
)

type Repository interface {
	Save(ctx context.Context, page *page.Page) error
	Init(ctx context.Context) error
	Show(ctx context.Context, page *page.Page) ([]page.Page, error)
}
