package main

import (
	"context"
	"fmt"
	"log"
	"time"

	gcron "github.com/go-co-op/gocron"
	pb "github.com/seswarrajan/grpc-mesh-demo/proto"
	"google.golang.org/grpc"
)

var (
	client pb.PaymentsServiceClient
	ctx    context.Context
	cancel context.CancelFunc
	count  float64
)

func main() {
	// In-cluster you'd use the service DNS 'payments:50051'
	target := "payments:50051"
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client = pb.NewPaymentsServiceClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	schedule := gcron.NewScheduler(time.UTC)
	schedule.SingletonModeAll()

	_, err = schedule.Every("1m").Do(func() {
		count = count + 1
		resp, err := client.ProcessPayment(ctx, &pb.PaymentRequest{
			OrderId:  "ORD-" + fmt.Sprintf("%f", count),
			Amount:   count,
			Currency: "USD",
		})
		if err != nil {
			log.Fatalf("ProcessPayment error: %v", err)
		}
		log.Printf("Payment result: status=%s txn=%s", resp.Status, resp.TransactionId)
	})
	if err != nil {
		log.Printf("\nerror in starting scheduler :%v\n", err)
		return
	}

	schedule.StartBlocking()

	select {} /// Forever loop
}
