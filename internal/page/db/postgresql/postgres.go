package postgres

import (
	"context"
	"fmt"
	"projects/DAB/internal/page"
	pageDB "projects/DAB/internal/page/storage"
	"projects/DAB/pkg/logging"
	"projects/DAB/pkg/storage/postgres"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

// Save сохраняет ивент в бд.
func (r *repository) Save(ctx context.Context, p *page.Page) error {
	q := `INSERT INTO tgbot_db (user_name, url, description, category, price, time_duration) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))

	if err := r.client.QueryRow(ctx, q, p.UserName, p.URL, p.Description, p.Category, p.Price, p.TimeDuration).Scan(&p.ID); err != nil {
		r.logger.Error(err)
		return err
	}

	return nil

}

// Show подбирает подходящие ивенты из бд.
func (r *repository) Show(ctx context.Context, p *page.Page) ([]page.Page, error) {
	q := `SELECT url, user_name, create_time FROM public.tgbot_db 
	WHERE category = $1
	AND price = $2
	AND time_duration = $3`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))

	rows, err := r.client.Query(ctx, q, p.Category, p.Price, p.TimeDuration)
	if err != nil {
		return nil, postgres.ErrNoSavedPages
	}

	pages := make([]page.Page, 0)

	for rows.Next() {
		var p page.Page

		err = rows.Scan(&p.URL, &p.UserName, &p.CreateTime)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
		r.logger.Info(p)
	}

	if len(pages) == 0 {
		return nil, postgres.ErrNoSavedPages
	}

	return pages, nil
}

func (r *repository) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS public.tgbot_db (id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_time DATE DEFAULT CURRENT_DATE,
    user_name VARCHAR(255),
    url VARCHAR(255),
	description VARCHAR(255),
	category VARCHAR(255),
	price INT,
	time_duration INT)`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))

	init, err := r.client.Exec(ctx, q)
	if err != nil {
		return err
	}
	r.logger.Info(init)

	return nil
}

func New(client postgres.Client, logger *logging.Logger) pageDB.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
