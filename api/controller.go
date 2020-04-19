package api

import (
	"net/http"

	"github.com/n4wei/memo/api/memo"
	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
)

type controller struct {
	dbClient db.Client
	router   *http.ServeMux
	logger   logger.Logger
}

func NewController(dbClient db.Client, logger logger.Logger) http.Handler {
	memoHandler := memo.NewHandler(dbClient, logger)

	router := http.NewServeMux()
	router.Handle("/memo", memoHandler)
	router.Handle("/memo/", memoHandler)

	return &controller{
		dbClient: dbClient,
		router:   router,
		logger:   logger,
	}
}

func (this *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.router.ServeHTTP(w, r)
}
