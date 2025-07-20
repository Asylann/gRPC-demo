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
	if err != nil {
		log.Fatal(err.Error())
	}

	c := pb.NewCartServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r1, err := c.AddToCart(ctx, &pb.AddToCardRequest{Item: &pb.Cart_Item{CartId: 2, Product: &pb.Product{
		Id:          2,
		Name:        "Sdjn",
		Description: "d,mkd",
		Price:       123.43,
		Size:        2,
		CategoryId:  2,
		ImageUrl:    "aaserfmks",
		SellerId:    4,
	}}})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(r1.IsAdded)
}
