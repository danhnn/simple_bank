package api

import (
	db "github.com/danhnn/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serve http request for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new http server and serve routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// add routers to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
