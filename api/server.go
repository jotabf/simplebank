package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jotabf/simplebank/db/sqlc"
	"github.com/jotabf/simplebank/token"
	"github.com/jotabf/simplebank/util"
)

// Server: serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer: creates a new HTTP server and setup the routing
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("user", validUser)
	}

	server.setupeRouters()

	return server, nil
}

func (server *Server) setupeRouters() {

	router := gin.Default()
	authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.GET("/users", server.getUser)

	router.POST("/tokens/renew", server.renewAccessToken)

	authRouters.POST("/accounts", server.createAccount)
	authRouters.GET("/accounts/:id", server.getAccount)
	authRouters.GET("/accounts", server.listAccount)
	// authRouters.DELETE("/accounts/:id", server.deleteAccount)

	authRouters.POST("/transfers", server.createTranfer)

	server.router = router
}

// Start: runs the HTTP server on a especific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
