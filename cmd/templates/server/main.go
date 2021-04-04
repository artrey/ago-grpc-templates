package main

import (
	"context"
	"github.com/artrey/ago-grpc-templates/cmd/templates/server/app"
	templatesV1Pb "github.com/artrey/ago-grpc-templates/pkg/api/proto/v1"
	"github.com/artrey/ago-grpc-templates/pkg/templates"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "9999"
	defaultDsn  = "postgres://localhost:5432/db"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dsn = defaultDsn
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr, dsn string) error {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	templatesSvc := templates.New(conn)
	server := app.NewServer(templatesSvc)
	grpcServer := grpc.NewServer()
	templatesV1Pb.RegisterTemplatesServiceServer(grpcServer, server)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return grpcServer.Serve(listener)
}
