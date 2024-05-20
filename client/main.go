package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/viictormg/grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)

	DoBidirectionalStreaming(c)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "1",
	}
	res, err := c.GetTest(context.Background(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response from server", res)
}

func DoClienteStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q20",
			Answer:   "Viictor",
			Question: "What is your name?",
			TestId:   "1",
		},
		{
			Id:       "q21",
			Answer:   "30",
			Question: "How old are you?",
			TestId:   "1",
		},
	}

	stream, err := c.SetQuestions(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for _, question := range questions {
		fmt.Println("Sending question: ", question.Id)
		stream.Send(question)
		time.Sleep(2 * time.Second)
	}

	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response from server", msg)
}

func DoServerStreaming(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "1",
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)

	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Response from server", msg)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {
	test := testpb.TakeTestRequest{Answer: "hola"}
	numerOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for i := 0; i < numerOfQuestions; i++ {
			err = stream.Send(&test)
			time.Sleep(2 * time.Second)

			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal("Error reading", err)
			}

			fmt.Println("Response from server", msg)
		}

		close(waitChannel)
	}()

	fmt.Println("FINALIZO")
	<-waitChannel
}
