package main

import (
	"context"
	"log"
	"time"

	cron "github.com/go-co-op/gocron/v2"
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

	// create a scheduler
	s, err := cron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.NewJob(
		cron.DurationJob(15*time.Second),
		cron.NewTask(func() {
			resp, err := client.ProcessPayment(ctx, &pb.PaymentRequest{
				OrderId:  "ORDER-001",
				Amount:   100,
				Currency: "USD",
			})
			if err != nil {
				log.Printf("[ERR] ProcessPayment error: %v", err)
			}
			log.Printf("Payment result: status=%s txn=%s", resp.Status, resp.TransactionId)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()

	select {} // wait forever
}
