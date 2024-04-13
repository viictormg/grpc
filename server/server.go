package server

import (
	"context"

	"github.com/viictormg/grpc/models"
	"github.com/viictormg/grpc/repository"
	"github.com/viictormg/grpc/studentpb"
)

type Server struct {
	repo repository.Repository
	studentpb.UnimplementedStudentServiceServer
}

func NewStudentServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {

	student, err := s.repo.GetStudent(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		},
		nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {

	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	err := s.repo.SetStudent(ctx, student)

	if err != nil {
		return nil, err
	}

	return &studentpb.SetStudentResponse{Id: req.Id}, nil
}
