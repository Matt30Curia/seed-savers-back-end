package order

import (
	"backend/seed-savers/services/auth"
	"backend/seed-savers/types"
	"backend/seed-savers/utils"
	"fmt"
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	usersStore   types.UserStore
	seedStore    types.SeedStore
	sessionStore *auth.AuthStore
}

func NewHandler(s types.OrderStore, us types.UserStore, seedStore types.SeedStore, sessionStore *auth.AuthStore) *Handler {
	return &Handler{s, us, seedStore, sessionStore}
}

func (h *Handler) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/create-order", auth.WithJWTAuth(h.handleCreateOrder, h.usersStore, h.sessionStore)).Methods("POST")
	router.HandleFunc("/update-order", auth.WithJWTAuth(h.handleUpdateOrder, h.usersStore, h.sessionStore)).Methods("PUT")
	router.HandleFunc("/orders-to-ship", auth.WithJWTAuth(h.handleOrdersToShip, h.usersStore, h.sessionStore)).Methods("GET")
	router.HandleFunc("/orders-to-recive", auth.WithJWTAuth(h.handleOrdersToRecive, h.usersStore, h.sessionStore)).Methods("GET")
	router.HandleFunc("/orders-delete/{id}", auth.WithJWTAuth(h.handleOrdersDelete, h.usersStore, h.sessionStore)).Methods("DELETE")
}

func (h *Handler) handleOrdersDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.store.DeleteOrder(id); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

}

func (h *Handler) handleOrdersToShip(w http.ResponseWriter, r *http.Request) {

	sender, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	orders, err := h.store.GetOrdersToBeSent(sender)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, orders)
}

func (h *Handler) handleOrdersToRecive(w http.ResponseWriter, r *http.Request) {
	reciver, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	orders, err := h.store.GetIncomingOrders(reciver)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, orders)
}

func (h *Handler) handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	var order types.Order

	payload, err := utils.DecodePayload[types.UpdateOrderPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	order = types.Order{ID: payload.OrderId, State: payload.State, Seed: types.Seed{Quantity: payload.SeedQuantity}}

	err = h.store.ModifyOrder(&order)
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
	
	quantity := h.seedStore.UserSeedQuantity(payload.SenderID, payload.SeedID)

	//verifico che l'utente che dovra spedire semi ha più semi ddi quanti richiesti TO-DO: gestire i crediti
	if quantity >= payload.SeedQuantity {
		err = h.store.MakeOrder(payload.SenderID, reciver, payload.SeedID, payload.SeedQuantity)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		
		err = h.usersStore.ModifySeedQuantity(&types.Seed{ID: payload.SeedID, Quantity: quantity - payload.SeedQuantity}, payload.SenderID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("non ci sono abbastanza semi"))
		return
	}

}
