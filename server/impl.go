package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"time"
	pb "todo_app/proto/v2"
	"todo_app/server/helpers"
)

func (s *server) AddTask(_ context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	if len(in.Description) == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"expected a task description, got an empty string",
		)
	}
	if in.DueDate.AsTime().Before(time.Now().UTC()) {
		return nil, status.Error(
			codes.InvalidArgument,
			"expected a task due_date that is in the future",
		)
	}
	id, err := s.d.addTask(in.Description, in.DueDate.AsTime())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"unexpected error: %s",
			err.Error(),
		)
	}
	return &pb.AddTaskResponse{Id: id}, nil
}

func (s *server) ListTasks(req *pb.ListTasksRequest, stream pb.TodoService_ListTasksServer) error {
	ctx := stream.Context()
	return s.d.getTasks(func(t any) error {
		// check for cancel/deadline
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				log.Printf("request canceled: %s", ctx.Err())
			case context.DeadlineExceeded:
				log.Printf("request deadline exceeded: %s", ctx.Err())
			default:
				log.Printf("request cancelled %s", ctx.Err())
			}
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}
		task := t.(*pb.Task)

		fmt.Printf(task.String())
		// use the filter for field mask
		helpers.Filter(task, req.Mask)

		overdue := task.DueDate != nil && !task.Done &&
			task.DueDate.AsTime().Before(time.Now().UTC())
		err := stream.Send(&pb.ListTasksResponse{
			Task:    task,
			Overdue: overdue,
		})
		return err
	})
}

func (s *server) UpdateTasks(stream pb.TodoService_UpdateTasksServer) error {
	totalLength := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Total: ", totalLength)
			return stream.SendAndClose(&pb.UpdateTasksResponse{})
		}
		if err != nil {
			return err
		}
		// to get the length of the data that is  transported
		// in the protobuf format
		out, _ := proto.Marshal(req)
		totalLength += len(out)
		s.d.updateTask(
			req.Id,
			req.Description,
			req.DueDate.AsTime(),
			req.Done,
		)
	}
}

func (s *server) DeleteTasks(stream pb.TodoService_DeleteTasksServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = s.d.deleteTask(req.Id)
		if err != nil {
			stream.Send(&pb.DeleteTasksResponse{
				Id:      req.Id,
				Success: false,
				Error:   err.Error(),
			})
		}
		stream.Send(&pb.DeleteTasksResponse{Id: req.Id, Success: true, Error: ""})
	}
}
