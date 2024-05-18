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

func (p *PostgresRepository) SetTest(ctx context.Context, test *models.Test) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO tests (id, name) VALUES ($1, $2)",
		test.Id, test.Name)

	return err
}

func (p *PostgresRepository) GetTest(ctx context.Context, id string) (*models.Test, error) {
	var test = models.Test{}

	row, err := p.db.QueryContext(ctx, "SELECT id, name FROM tests WHERE id = $1", id)
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
		err := row.Scan(&test.Id, &test.Name)
		if err != nil {
			return nil, err
		}
	}

	return &test, nil
}

func (p *PostgresRepository) SetQuestion(ctx context.Context, question *models.Question) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO questions (id, answer, question, test_id) VALUES ($1, $2, $3, $4)",
		question.Id, question.Answer, question.Question, question.TestId)

	return err
}

func (p *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error) {

	rows, err := p.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id IN (SELECT student_id FROM enrollments WHERE test_id = $1)", testId)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	var students []*models.Student
	for rows.Next() {
		var student = models.Student{}
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}
	return students, nil
}

func (p *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO enrollments (student_id, test_id) VALUES ($1, $2)",
		enrollment.StudentId, enrollment.TestId)

	return err
}

func (p *PostgresRepository) GetQuestionsPerTest(ctx context.Context, testId string) ([]*models.Question, error) {

	rows, err := p.db.QueryContext(ctx, "SELECT id, question FROM questions WHERE test_id = $1", testId)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	var questions []*models.Question
	for rows.Next() {
		var question = models.Question{}
		err := rows.Scan(&question.Id, &question.Question)
		if err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}
	return questions, nil
}
