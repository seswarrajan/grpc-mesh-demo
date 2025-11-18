package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/seswarrajan/grpc-mesh-demo/proto"
	"google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var DeploymentLabel string

type server struct {
	pb.UnimplementedPaymentsServiceServer
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("Processing payment for order %s, amount %.2f %s", req.OrderId, req.Amount, req.Currency)
	return &pb.PaymentResponse{
		Status:          "SUCCESS",
		TransactionId:   fmt.Sprintf("txn-%s", req.OrderId),
		DeploymentLabel: DeploymentLabel,
	}, nil
}

func main() {
	port := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentsServiceServer(grpcServer, &server{})

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// log.Println("Serving gRPC on server:v1")
	// log.Println("Serving gRPC on server:v2")

	kubeconfig := flag.String("kubeconfig", "~/.kube/config", "path to kubeconfig")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deployClient := clientset.AppsV1().Deployments("payments")

	deploy, err := deployClient.Get(context.TODO(), "payments-v2", v1.GetOptions{})
	if err != nil {
		panic(err)
	}

	labels := deploy.Labels
	for k, v := range labels {
		DeploymentLabel += fmt.Sprintf("%s=%s,", k, v)
	}

	// remove trailing comma
	if len(DeploymentLabel) > 0 {
		DeploymentLabel = DeploymentLabel[:len(DeploymentLabel)-1]
	}

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
