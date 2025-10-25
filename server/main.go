package main

import (
	"google.golang.org/grpc"
	pb "todo_app/todo/v1"
)

func main() {
	s := grpc.NewServer(opts...)
	pb.RegisterTodoServiceServer(s, &server{
		d: New(),
	})
	defer s.Stop()
}
