package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jotabf/simplebank/api"
	db "github.com/jotabf/simplebank/db/sqlc"
	"github.com/jotabf/simplebank/gapi"
	"github.com/jotabf/simplebank/pb"
	"github.com/jotabf/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := db.NewStore(connDB)
	// runGinServer(config, store)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config *util.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listner, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal("Cannot create listener: ", err)
	}

	log.Printf("Starting grpc server at %s", listner.Addr().String())
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Cannot start gRPC server: ", err)
	}
}

func runGatewayServer(config *util.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Cannot register simple bank handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger/"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listner, err := net.Listen("tcp", config.HTTPServerAddr)
	if err != nil {
		log.Fatal("Cannot create listener: ", err)
	}

	log.Printf("Starting HTTP gateway server at %s", listner.Addr().String())
	err = http.Serve(listner, mux)
	if err != nil {
		log.Fatal("Cannot start HTTP gateway server: ", err)
	}
}

func runGinServer(config *util.Config, store db.Store) {

	server, err := api.NewServer(*config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddr)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
