package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
	pb "todo_app/proto/v2"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: client [IP_ADDR]")
	}

	// make a grpc client connection body
	addr := args[0]
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}(conn)

	// make a client for TodoServiceClient
	c := pb.NewTodoServiceClient(conn)
	// add a task
	cmd := args[1]
	switch cmd {
	case "add":
		fmt.Println("...........ADD............")
		dueDate := time.Now().Add(5 * time.Second)
		id := addTask(c, "This is a task", dueDate)
		fmt.Println("Added Task:", id)
		fmt.Println("..........................")
	case "print":
		fmt.Println("...........LIST...........")
		printTasks(c)
		fmt.Println("..........................")
	case "update":
		fmt.Println("...........UPDATE.........")
		updateTasks(c, &pb.UpdateTasksRequest{
			Id: 1, Description: "oombikko myre",
			DueDate: tpb.New(time.Now().Add(3 * time.Hour)),
			Done:    false,
		}, &pb.UpdateTasksRequest{
			Id: 2, Description: "oombikko myre 2 times",
			DueDate: tpb.New(time.Now().Add(1 * time.Hour)),
		}, &pb.UpdateTasksRequest{
			Id: 3, Description: "oombikko 3myre",
		})
		fmt.Println("..........................")
	case "delete":
		fmt.Println("...........DELETE.........")
		deleteTasks(c,
			&pb.DeleteTasksRequest{Id: 1},
			&pb.DeleteTasksRequest{Id: 2},
			&pb.DeleteTasksRequest{Id: 3},
		)
		fmt.Println("..........................")

	}

}

func addTask(c pb.TodoServiceClient, description string, dueDate time.Time) uint64 {
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     tpb.New(dueDate),
	}
	res, err := c.AddTask(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return res.Id
}

func printTasks(c pb.TodoServiceClient) {
	req := &pb.ListTasksRequest{}
	stream, err := c.ListTasks(context.Background(), req)

	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("unexpted error: %v", err)
		}
		fmt.Println(res.Task.String(), "overdue: ", res.Overdue)
	}
}

func updateTasks(c pb.TodoServiceClient, reqs ...*pb.UpdateTasksRequest) {
	stream, err := c.UpdateTasks(context.Background())
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	for _, req := range reqs {
		err := stream.Send(req)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
		if req != nil {
			fmt.Printf("updated task with id: %d\n", req.Id)
		}
	}
	if _, err = stream.CloseAndRecv(); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
}

func deleteTasks(c pb.TodoServiceClient, reqs ...*pb.DeleteTasksRequest) {
	stream, err := c.DeleteTasks(context.Background())
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v\n", err)
			}
			log.Println("deleted task")
			fmt.Println(res)
		}
	}()
	for _, req := range reqs {
		if err := stream.Send(req); err != nil {
			return
		}
	}

	if err := stream.CloseSend(); err != nil {
		return
	}
	<-waitc
}
