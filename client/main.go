package main

import (
	"context"
	"log"
	"time"

	pb "github.com/seswarrajan/grpc-mesh-demo/proto"
	"google.golang.org/grpc"
)

func main() {
	// In-cluster you'd use the service DNS 'payments:50051'
	target := "payments:50051"
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.PaymentRequest{
		OrderId:  "ORD-1001",
		Amount:   199.99,
		Currency: "USD",
	}

	resp, err := client.ProcessPayment(ctx, req)
	if err != nil {
		log.Fatalf("ProcessPayment error: %v", err)
	}
	log.Printf("Payment result: status=%s txn=%s", resp.Status, resp.TransactionId)
}
