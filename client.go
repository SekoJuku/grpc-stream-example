package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	pb "github.com/SekoJuku/grpc-stream-example/proto/message"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	c := pb.NewServiceClient(conn)

	SendMessage(c)
}

func SendMessage(client pb.ServiceClient) error {
	fmt.Println("Client is starting!")
	stream, _ := client.SendMessage(context.Background())
	defer stream.CloseSend()

	scanner := bufio.NewScanner(os.Stdin)
	// scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text := scanner.Text()
		err := stream.Send(&pb.Message{Text: text})
		if err != nil {
			return err
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println("Streaming Recv:", res.Text)
	}
	return nil
}
