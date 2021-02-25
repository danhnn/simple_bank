package api

import (
	"os"
	"testing"

	db "github.com/danhnn/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	// config := util.Config{
	// 	TokenSymmetricKey:   util.RandomString(32),
	// 	AccessTokenDuration: time.Minute,
	// }

	server := NewServer(store)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
