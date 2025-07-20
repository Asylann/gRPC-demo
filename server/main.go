package main

import (
	"context"
	pb "github.com/Asylann/gRPC_Demo/proto"
	"github.com/Asylann/gRPC_Demo/server/internal/config"
	"github.com/Asylann/gRPC_Demo/server/internal/models"
	"github.com/Asylann/gRPC_Demo/server/internal/repository"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	pb.UnimplementedCartServiceServer
	CartStore *repository.CartStore
}

func (s *Server) CreateCart(ctx context.Context, req *pb.CreateCartRequest) (*pb.CreateCartResponse, error) {
	err := s.CartStore.CreateCart(int(req.GetUserId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.CreateCartResponse{IsCreated: true}, nil
}

func (s *Server) AddToCart(ctx context.Context, req *pb.AddToCardRequest) (*pb.AddToCardResponse, error) {
	product := models.Product{
		ID:          int(req.GetItem().Product.Id),
		Name:        req.GetItem().Product.Name,
		Description: req.GetItem().Product.Description,
		Price:       req.GetItem().Product.Price,
		Size:        int(req.GetItem().Product.Size),
		CategoryID:  int(req.GetItem().Product.CategoryId),
		ImageURL:    req.GetItem().Product.ImageUrl,
		SellerID:    int(req.GetItem().Product.SellerId),
	}
	log.Println(product)
	err := s.CartStore.AddToCart(int(req.GetItem().CartId), product)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.AddToCardResponse{IsAdded: true}, nil
}

func (s *Server) GetCartById(ctx context.Context, req *pb.GetCartByIdRequest) (*pb.GetCartByIdResponse, error) {
	cart, err := s.CartStore.GetCartById(int(req.GetId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.GetCartByIdResponse{Cart: &pb.Cart{Id: int32(cart.Id), UserId: int32(cart.User_id)}}, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	db := repository.InitDB(cfg)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	cartStore, err := repository.NewCartStore(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	srv := &Server{
		CartStore: &cartStore,
	}

	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, srv)

	quit, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Println("Server is running on :50051")
		if err = s.Serve(l); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-quit.Done()
	log.Println("Shut down processing...")

	done := make(chan interface{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Server is stopped running !")
	case <-time.After(10 * time.Second):
		log.Println("Server is stopped running due to timeout !")
		s.Stop()
	}

}
