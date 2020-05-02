package api

import (
	"net/http"

	"github.com/n4wei/memo/api/handler"
	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
)

type controller struct {
	dbClient db.Client
	router   http.Handler
	logger   logger.Logger
}

func New(dbClient db.Client, logger logger.Logger) http.Handler {
	userHandler := handler.NewUserHandler(dbClient, logger)
	memoHandler := handler.NewMemoHandler(dbClient, logger)

	router := newRouter()
	router.addRoute("/user", userHandler)
	router.addRoute("/user/:guid/memo", memoHandler)

	return &controller{
		dbClient: dbClient,
		router:   router,
		logger:   logger,
	}
}

func (this *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.router.ServeHTTP(w, r)
}
