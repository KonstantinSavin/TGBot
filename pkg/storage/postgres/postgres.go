package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"projects/DAB/internal/config"
	"projects/DAB/pkg/utils"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

var ErrNoSavedPages = errors.New("нет сохраненных ссылок")

func New(ctx context.Context, maxAttempts int, sc config.StorageConfig) (*pgxpool.Pool, error) {

	var pool *pgxpool.Pool
	var err error

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("ошибка с попытками подключиться к postgresql")
	}

	return pool, nil
}
