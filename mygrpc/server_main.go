package main

import (
        "log"
        "net"

        "golang.org/x/net/context"
        "google.golang.org/grpc"
        pb "./pb"
        "google.golang.org/grpc/reflection"
        "google.golang.org/grpc/credentials"
)

const (
        port = ":8000"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
        return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
        return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}

func main() {
        creds, _ := credentials.NewServerTLSFromFile("/home/di_sun/CA/grpc-server.pem", "/home/di_sun/CA/grpc-server-key.pem")
        lis, err := net.Listen("tcp", port)
        if err != nil {
                log.Fatalf("failed to listen: %v", err)
        }
        s := grpc.NewServer(grpc.Creds(creds))
        pb.RegisterGreeterServer(s, &server{})
        // Register reflection service on gRPC server.
        reflection.Register(s)
        if err := s.Serve(lis); err != nil {
                log.Fatalf("failed to serve: %v", err)
        }
}
