package main

import (
	tpb "google.golang.org/protobuf/types/known/timestamppb"
	"time"
	pb "todo_app/todo/v1"
)

type inMemoryDb struct {
	tasks []*pb.Task
}

func New() db {
	return &inMemoryDb{}
}

func (d *inMemoryDb) addTask(description string, dueDate time.Time) (uint64, error) {
	nextId := uint64(len(d.tasks) + 1)
	task := &pb.Task{
		Id:          nextId,
		Description: description,
		DueDate:     tpb.New(dueDate),
	}
	d.tasks = append(d.tasks, task)
	return nextId, nil
}
