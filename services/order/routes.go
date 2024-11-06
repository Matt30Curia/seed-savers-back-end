package order

import (
	"backend/seed-savers/services/auth"
	//"backend/seed-savers/services/order"
	"backend/seed-savers/types"
	"backend/seed-savers/utils"

	"log"

	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	usersStore   types.UserStore
	sessionStore *auth.AuthStore
}

func NewHandler(s types.OrderStore, us types.UserStore, sessionStore *auth.AuthStore) *Handler {
	return &Handler{s, us, sessionStore}
}

func (h *Handler) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/create-order", auth.WithJWTAuth(h.handleCreateOrder, h.usersStore, h.sessionStore)).Methods("POST")
	router.HandleFunc("/update-order", auth.WithJWTAuth(h.handleModifyOrder, h.usersStore, h.sessionStore)).Methods("PUT")
	router.HandleFunc("/sender-order", auth.WithJWTAuth(h.handleSenderOrder, h.usersStore, h.sessionStore)).Methods("GET")
	router.HandleFunc("/reciver-order", auth.WithJWTAuth(h.handleReciverOrder, h.usersStore, h.sessionStore)).Methods("GET")
}

func (h *Handler) handleReciverOrder(w http.ResponseWriter, r *http.Request) {

	sender, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	orders, err := h.store.GetOrdersBySender(sender)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, orders)
}

func (h *Handler) handleSenderOrder(w http.ResponseWriter, r *http.Request) {
	reciver, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	orders, err := h.store.GetOrdersByReciver(reciver)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, orders)
}

func (h *Handler) handleModifyOrder(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.DecodePayload[types.Order](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	reciver, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	payload.Reciver.ID = reciver

	err = h.store.ModifyOrder(payload)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.DecodePayload[types.OrderPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	reciver, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	err = h.store.MakeOrder(reciver, payload.SenderID, payload.SeedID, payload.SeedQuantity)
	if err != nil {
		log.Println(err)
	}
}
