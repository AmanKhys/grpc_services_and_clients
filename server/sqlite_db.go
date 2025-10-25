package main

//
// import (
// 	"database/sql"
// 	_ "modernc.org/sqlite"
// )
//
// import (
// 	tpb "google.golang.org/protobuf/types/known/timestamppb"
// 	"time"
// 	pb "todo_app/todo/v1"
// )
//
// type sqliteDB struct {
// 	db sql.DB
// }
//
// func NewSqliteDB() sqliteDB {
// }
//
// func (d *sqliteDB) addTask(description string, dueDate time.Time) (uint64, error) {
// 	nextId := uint64(len(d.tasks) + 1)
// 	task := &pb.Task{
// 		Id:          nextId,
// 		Description: description,
// 		DueDate:     tpb.New(dueDate),
// 	}
// 	d.tasks = append(d.tasks, task)
// 	return nextId, nil
// }
