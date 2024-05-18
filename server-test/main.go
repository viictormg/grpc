package main

import (
	"fmt"
	"log"
	"net"

	"github.com/viictormg/grpc/database"
	"github.com/viictormg/grpc/server"
	"github.com/viictormg/grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgresql:tx7L8AMk91SDNS@localhost:5432/prenlink-db?sslmode=disable")

	fmt.Println("SERVER TEST RUNNING ON PORT 5070")
	if err != nil {
		log.Fatal(err)
	}
	srv := server.NewTestServer(repo)

	s := grpc.NewServer()

	testpb.RegisterTestServiceServer(s, srv)
	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}

}
