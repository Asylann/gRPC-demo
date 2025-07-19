package main

import (
	"context"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.UnimplementedCalculatorServer
}

func (s *Server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	sum := req.GetNum1() + req.GetNum2()
	return &pb.AddResponse{Sum: sum}, nil
}

func (s *Server) Subtract(ctx context.Context, req *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	difference := req.GetNum1() - req.GetNum2()
	return &pb.SubtractResponse{Difference: difference}, nil
}

func main() {
	s := grpc.NewServer()
	srv := &Server{}
	pb.RegisterCalculatorServer(s, srv)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Server is running on :50051")
	if err = s.Serve(l); err != nil {
		log.Fatal(err.Error())
	}

}
