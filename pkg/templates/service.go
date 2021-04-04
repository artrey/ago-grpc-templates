package templates

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

type Template struct {
	Id      int64
	Title   string
	Phone   string
	Created int64
	Updated int64
}

func New(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s Service) Create(ctx context.Context, title, phone string) (*Template, error) {
	template := Template{
		Title: title,
		Phone: phone,
	}
	err := s.pool.QueryRow(ctx, `
INSERT INTO templates(title, phone) VALUES ($1, $2)
RETURNING id, EXTRACT(epoch FROM created)::integer, EXTRACT(epoch FROM updated)::integer
	`, title, phone).Scan(&template.Id, &template.Created, &template.Updated)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &template, nil
}

func (s Service) List(ctx context.Context) ([]*Template, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, title, phone, EXTRACT(epoch FROM created)::integer, EXTRACT(epoch FROM updated)::integer
FROM templates LIMIT 100`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	templates := make([]*Template, 0)
	for rows.Next() {
		var t Template
		err = rows.Scan(&t.Id, &t.Title, &t.Phone, &t.Created, &t.Updated)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		templates = append(templates, &t)
	}

	return templates, nil
}

func (s Service) Get(ctx context.Context, id int64) (*Template, error) {
	template := Template{
		Id: id,
	}
	err := s.pool.QueryRow(ctx, `
SELECT title, phone, EXTRACT(epoch FROM created)::integer, EXTRACT(epoch FROM updated)::integer
FROM templates WHERE id=$1
	`, id).Scan(&template.Title, &template.Phone, &template.Created, &template.Updated)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &template, nil
}

func (s Service) Update(ctx context.Context, id int64, title, phone string) (*Template, error) {
	template := Template{
		Id:    id,
		Title: title,
		Phone: phone,
	}
	err := s.pool.QueryRow(ctx, `
UPDATE templates SET title=$1, phone=$2 WHERE id=$3
RETURNING EXTRACT(epoch FROM created)::integer, EXTRACT(epoch FROM updated)::integer
	`, title, phone, id).Scan(&template.Created, &template.Updated)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &template, nil
}

func (s Service) Delete(ctx context.Context, id int64) (*Template, error) {
	template := Template{
		Id: id,
	}
	err := s.pool.QueryRow(ctx, `
DELETE FROM templates WHERE id=$1
RETURNING title, phone, EXTRACT(epoch FROM created)::integer, EXTRACT(epoch FROM updated)::integer
	`, id).Scan(&template.Title, &template.Phone, &template.Created, &template.Updated)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &template, nil
}
