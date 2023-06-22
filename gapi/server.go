package gapi

import (
	"fmt"

	db "github.com/jotabf/simplebank/db/sqlc"
	"github.com/jotabf/simplebank/pb"
	"github.com/jotabf/simplebank/token"
	"github.com/jotabf/simplebank/util"
)

// Server: serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer: creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
