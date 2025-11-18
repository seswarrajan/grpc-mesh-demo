package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/seswarrajan/grpc-mesh-demo/proto"
	"google.golang.org/grpc"
)

var (
	client pb.PaymentsServiceClient
	count  float64
)

func main() {
	// In-cluster you'd use the service DNS 'payments:50051'
	target := "payments.payments.svc.cluster.local:50051"
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client = pb.NewPaymentsServiceClient(conn)

	for {
		count = count + 1
		resp, err := client.ProcessPayment(context.Background(), &pb.PaymentRequest{
			OrderId:  "ORD-" + fmt.Sprintf("%f", count),
			Amount:   count,
			Currency: "USD",
		})
		if err != nil {
			log.Fatalf("ProcessPayment error: %v", err)
		}
		log.Printf("Payment result: status=%s txn=%s Processed_By_Deployment_Label=%v", resp.Status, resp.TransactionId, resp.DeploymentLabel)
		time.Sleep(10 * time.Second)
	}
}
