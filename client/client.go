package main

import (
	"context"
	"flag"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	if flag.NArg() < 3 {
		log.Fatal("You need to enter two numbers!")
	}

	operation := flag.Arg(0)

	num1, err := strconv.ParseFloat(flag.Arg(1), 64)
	if err != nil {
		log.Fatal(err.Error())
	}

	num2, err := strconv.ParseFloat(flag.Arg(2), 64)
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	c := pb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if operation == "Addition" {
		r, err := c.Add(ctx, &pb.AddRequest{Num1: num1, Num2: num2})
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("The sum of %v and %v is %v", num1, num2, r.GetSum())
	} else if operation == "Subtraction" {
		r, err := c.Subtract(ctx, &pb.SubtractRequest{Num1: num1, Num2: num2})
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("The diffrence between %v and %v is %v", num1, num2, r.GetDifference())
	}
}
