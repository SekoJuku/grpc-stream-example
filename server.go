package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/SekoJuku/grpc-stream-example/proto/message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var activeConnections map[string]pb.Service_SendMessageServer

type server struct {
	pb.UnimplementedServiceServer
}

type streamManager struct {
}

func main() {
	activeConnections = make(map[string]pb.Service_SendMessageServer)
	fmt.Println("Server is starting!")

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterServiceServer(s, &server{})
	s.Serve(listen)
}

func (s *server) SendMessage(stream pb.Service_SendMessageServer) error {

	activeConnections["sd"] = stream

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(res)
		stream.Send(&pb.Message{Text: res.Text})
	}
}
