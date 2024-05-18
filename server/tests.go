package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/viictormg/grpc/models"
	"github.com/viictormg/grpc/repository"
	"github.com/viictormg/grpc/studentpb"
	"github.com/viictormg/grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &testpb.Test{Id: test.Id, Name: test.Name}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.Id,
		Name: req.Name,
	}

	err := s.repo.SetTest(ctx, test)

	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{Id: req.Id, Name: req.Name}, nil
}

func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		question := &models.Question{
			Id:       msg.GetId(),
			Answer:   msg.GetAnswer(),
			Question: msg.GetQuestion(),
			TestId:   msg.GetTestId(),
		}

		err = s.repo.SetQuestion(context.Background(), question)

		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		erollment := &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}

		err = s.repo.SetEnrollment(context.Background(), erollment)

		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.TestId)

	if err != nil {
		return err
	}

	for _, s := range students {
		student := &studentpb.Student{
			Id:   s.Id,
			Name: s.Name,
			Age:  s.Age,
		}

		err := stream.Send(student)
		time.Sleep(2 * time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "1")

	if err != nil {
		return err
	}

	i := 0

	var currentQuestion = &models.Question{}

	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}

			err := stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		fmt.Println("Answer:", answer.GetAnswer())
	}
}
