package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/rahulsingh577/greetService-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("grpc Client started..")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to the server: %v", err)
	}

	defer conn.Close()
	client := greetpb.NewGreetServiceClient(conn)
	//	doUnary(client)
	//doServerStreaming(client)
	doCientStreaming(client)

}

func doUnary(client greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul",
			LastName:  "Singh",
		},
	}

	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling the Greet RPC: %v", err)
	}
	fmt.Printf("Response from greet: %f", res.Result)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul",
			LastName:  "Singh",
		},
	}

	res, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling the Greet many times RPC: %v",
			err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			//we have reached end ogf the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream : %v", err)
		}
		log.Printf("/n Response from GreetManyTimes : %v", msg.GetResult())
	}

}

func doCientStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming RPC..")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rahul",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tanuj",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anurag",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bojack",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aurthur",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Gatsby",
			},
		},
	}
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling longGreet: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while connecting to the longGreet: %v", err)
	}
	fmt.Printf("LongGreet Respone: %v\n", res)
}
