package api

import (
	"backend/seed-savers/services/auth"
	"backend/seed-savers/services/order"
	"backend/seed-savers/services/seed"

	"backend/seed-savers/services/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	adress string
	db     *sql.DB
}

func NewServer(adress string, db *sql.DB) *ApiServer {
	return &ApiServer{adress: adress, db: db}
}

func (a *ApiServer) Run() error {
	router := mux.NewRouter()

	authSessionStore := auth.NewOauth()
	userStore := user.NewStore(a.db)
	seedStore := seed.NewStore(a.db)
	orderStore := order.NewStore(a.db)

	userHandler := user.NewHandler(userStore, authSessionStore)
	seedHandler := seed.NewHandler(seedStore, userStore, authSessionStore)
	orderHandler := order.NewHandler(orderStore, userStore, authSessionStore)

	userHandler.RegisterRouter(router)
	seedHandler.RegisterRouter(router)
	orderHandler.RegisterRouter(router)

	log.Println("listening on: ", a.adress)
	return http.ListenAndServe(a.adress, router)
}
