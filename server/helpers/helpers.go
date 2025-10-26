package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"unsafe"
)

// serlalizedSize calculates the size in memory of data
// and the size after Protobuf serialization of the wrapper.
// It returns two values. The number of bytes for data in memory
// and the number of bytes after serialization of wrapper - 1
// (removes the byte for tag + wire type).
func serializedSize[D constraints.Integer, W protoreflect.ProtoMessage](data D, wrapper W) (uintptr, int) {
	out, err := proto.Marshal(wrapper)
	if err != nil {
		log.Fatal(err)
	}
	return unsafe.Sizeof(data), len(out) - 1
}

func main() {
	t := &pb.Tags{}
}
