package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/QuangNV23062004/grpc_example_golang/coffeeshop_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to gprc server, %s", err)
	}

	defer conn.Close()

	c := pb.NewCoffeeShopServiceClient(conn)
	log.Println("Connected to gRPC server successfully")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("Error while calling GetMenu: %s", err)
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			menuResp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Printf("Error while receiving menu item: %s", err)
				close(done)
				return
			}
			items = menuResp.Items
			log.Println("Received menu items:")
			for _, item := range items {
				log.Printf("Item ID: %s, Name: %s", item.GetId(), item.GetName())
			}
		}
	}()

	<-done

	receipt, err := c.PlaceOrder(ctx, &pb.OrderRequest{Items: items})
	if err != nil {
		log.Fatalf("Error while placing order: %s", err)
	}

	log.Printf("Order placed successfully, receipt ID: %s", receipt.GetId())

	orderStatus, err := c.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Fatalf("Error while getting order status: %s", err)
	}
	log.Printf("Order status for receipt ID %s: %s", orderStatus.GetOrderId(), orderStatus.GetStatus())
}
