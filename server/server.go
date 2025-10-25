package main

import (
	// "database/sql"
	"time"
	pb "todo_app/todo/v1"
)

type server struct {
	d db
	pb.UnimplementedTodoServiceServer
}

type db interface {
	addTask(description string, dueDate time.Time) (uint64, error)
	getTasks(func(any) error) error
	updateTask(id uint64, description string, dueDate time.Time, done bool) error
}
