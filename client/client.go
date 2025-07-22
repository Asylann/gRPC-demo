package main

import (
	"context"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	c := pb.NewCartServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r1, err := c.GetCartByUserId(ctx, &pb.GetCartByUserIdRequest{Id: int32(13)})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(r1.GetCart())
}
