package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/rahulsingh577/greetService-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

const port = "0.0.0.0:50051"

type server struct{}

func (s *server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("\n Greet unary func invoked.. : %v", in)
	firstName := in.GetGreeting().GetFirstName()
	result := "Hello" + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,
	stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Printf("\n Greet server stream func invoked.. : %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {

		result := " Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil

}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Printf("\n Long greet function was inovked with input : %v", stream)
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finsished reading client stream..

			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			//break
		}
		if err != nil {
			log.Printf("error while reading the strem : %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += " Hello " + firstName + " ! "
	}

	return nil
}

func main() {
	fmt.Print("Greet server stating.. ")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
