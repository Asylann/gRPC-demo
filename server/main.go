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
	log.Printf("Cart was created by user id = %v \n", req.UserId)
	return &pb.CreateCartResponse{IsCreated: true}, nil
}

func (s *Server) AddToCart(ctx context.Context, req *pb.AddToCardRequest) (*pb.AddToCardResponse, error) {
	product := models.Product{
		ID: int(req.GetItem().Product.Id),
	}
	log.Println(product)
	err := s.CartStore.AddToCart(int(req.GetItem().CartId), product)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.AddToCardResponse{IsAdded: true}, nil
}

func (s *Server) GetCartByUserId(ctx context.Context, req *pb.GetCartByUserIdRequest) (*pb.GetCartByUserIdResponse, error) {
	cart, err := s.CartStore.GetCartByUserId(int(req.GetId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.GetCartByUserIdResponse{Cart: &pb.Cart{Id: int32(cart.Id), UserId: int32(cart.User_id)}}, nil
}

func (s *Server) GetItemsOfCartById(ctx context.Context, req *pb.GetItemsOfCartByIdRequest) (*pb.GetItemsOfCartByIdResponse, error) {
	products, err := s.CartStore.GetProductsOfCartById(int(req.GetId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	resProducts := []*pb.Product{}
	for i := 0; i < len(products); i++ {
		resProducts = append(resProducts, &pb.Product{Id: int32(products[i])})
	}
	return &pb.GetItemsOfCartByIdResponse{Product: resProducts}, nil
}

func (s *Server) DeleteItemFromCart(ctx context.Context, req *pb.DeleteItemFromCartRequest) (*pb.DeleteItemFromCartResponse, error) {
	err := s.CartStore.DeleteItemFromCart(int(req.GetCartId()), int(req.GetProductId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.DeleteItemFromCartResponse{DeletedProduct: &pb.Product{Id: req.ProductId}}, nil
}

func (s *Server) GetEtagVersionByUserId(ctx context.Context, req *pb.GetEtagVersionByUserIdRequest) (*pb.GetEtagVersionByUserIdResponse, error) {
	version, err := s.CartStore.GetEtagVersionByUserId(int(req.GetUserId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.GetEtagVersionByUserIdResponse{Version: int64(version)}, nil
}

func (s *Server) ChangeEtagVersionByUserId(ctx context.Context, req *pb.ChangeEtagVersionByUserIdRequest) (*pb.ChangeEtagVersionByUserIdResponse, error) {
	changedVersion, err := s.CartStore.ChangeEtagVersionByUserId(int(req.GetUserId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.ChangeEtagVersionByUserIdResponse{IsChanged: true, ChangedVersion: int64(changedVersion)}, nil
}

func (s *Server) DeleteCart(ctx context.Context, req *pb.DeleteCartRequest) (*pb.DeleteCartResponse, error) {
	deletedCart, err := s.CartStore.DeleteCart(int(req.GetUserId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.DeleteCartResponse{DeletedCart: &pb.Cart{Id: int32(deletedCart.Id), UserId: int32(deletedCart.User_id)}}, nil
}

func (s *Server) DeleteProductOfCarts(ctx context.Context, req *pb.DeleteProductOfCartsRequest) (*pb.DeleteProductOfCartsResponse, error) {
	err := s.CartStore.DeleteProductOfCarts(int(req.GetProductId()))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.DeleteProductOfCartsResponse{ProductId: req.GetProductId()}, nil
}

func (s *Server) ChangeEtagVersionOfCartsByProductId(ctx context.Context, req *pb.ChangeEtagVersionOfCartsByProductIdRequest) (*pb.ChangeEtagVersionOfCartsByProductIdResponse, error) {
	err := s.CartStore.ChangeEtagVersionByProductId(int(req.ProductId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.ChangeEtagVersionOfCartsByProductIdResponse{IsChanged: true}, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	repository.InitDB(cfg)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	cartStore, err := repository.NewCartStore()
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
