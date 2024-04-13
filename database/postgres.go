package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/viictormg/grpc/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (p *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)",
		student.Id, student.Name, student.Age)

	return err
}

func (p *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	var student = models.Student{}
	row, err := p.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := row.Close()
		if err != nil {
			panic(err)
		}
	}()

	for row.Next() {
		err := row.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}
	}
	return &student, nil
}
