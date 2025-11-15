package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/example/grpc-mesh-demo/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPaymentsServiceServer
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("Processing payment for order %s, amount %.2f %s", req.OrderId, req.Amount, req.Currency)
	return &pb.PaymentResponse{
		Status:        "SUCCESS",
		TransactionId: fmt.Sprintf("txn-%s", req.OrderId),
	}, nil
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentsServiceServer(grpcServer, &server{})

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("gRPC server listening on %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc serve error: %v", err)
		}
	}()

	<-sigs
	log.Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
}
