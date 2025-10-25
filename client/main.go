package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
	pb "todo_app/todo/v1"
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
	fmt.Println("...........ADD............")
	dueDate := time.Now().Add(5 * time.Second)
	id := addTask(c, "This is a task", dueDate)
	fmt.Println("Added Task:", id)
	fmt.Println("..........................")

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
