package main

import (
	"time"
	pb "todo_app/todo/v1"
)

type server struct {
	d db
	pb.UnimplementedTodoServiceServer
}

type db interface {
	addTask(description string, dueDate time.Time) (uint64, error)
}
