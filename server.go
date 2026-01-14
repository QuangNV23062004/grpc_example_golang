package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/QuangNV23062004/grpc_example_golang/coffeeshop_proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServiceServer
}

func (s *server) GetMenu(req *pb.MenuRequest, res pb.CoffeeShopService_GetMenuServer) error {
	items := []*pb.Item{
		{Id: "1", Name: "Espresso"},
		{Id: "2", Name: "Cappuccino"},
		{Id: "3", Name: "Latte"},
	}

	log.Println("Sending menu items to client...")

	for i, _ := range items {
		res.Send(&pb.MenuResponse{Items: items[0 : i+1]})
	}

	return nil
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.OrderRequest) (*pb.Receipt, error) {
	return &pb.Receipt{Id: "12345"}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{OrderId: receipt.Id, Status: "In Progress"}, nil
}

func main() {
	//setup listener on 9001
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen on port 9001: %v", err)
	}

	fmt.Println("gRPC server listening on port 9001")
	grpcServer := grpc.NewServer()
	pb.RegisterCoffeeShopServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
