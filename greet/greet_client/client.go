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
	//doCientStreaming(client)
	doBidirectionalStreaming(client)

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

func doBidirectionalStreaming(client greetpb.GreetServiceClient) {

	fmt.Println("Starting bidriectional client streaming RPC..")

	waitc := make(chan struct{})

	// create a stream by invoking the client
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating the stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rahul",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tanuj",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anurag",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bojack",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aurthur",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Gatsby",
			},
		},
	}

	go func() {
		// we send messages to teh client
		for _, req := range requests {

			stream.Send(req)
			fmt.Printf("Sending message: %v\n", req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()

	}()

	// we do recieve from the clients

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
			}
			if err != nil {
				log.Fatalf("Error while recieving : %v", err)
				close(waitc)
			}
			fmt.Printf("Recived: %v", res.GetResult())

		}
		close(waitc)
	}()

	// block the things till processing is done
	<-waitc

}
