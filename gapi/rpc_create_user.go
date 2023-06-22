package gapi

import (
	"context"

	db "github.com/jotabf/simplebank/db/sqlc"
	"github.com/jotabf/simplebank/pb"
	"github.com/jotabf/simplebank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hash, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hash,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "Username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertToUserResponse(user),
	}

	return res, nil
}
