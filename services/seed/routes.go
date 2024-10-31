package seed

import (
	"backend/seed-savers/services/auth"
	"backend/seed-savers/types"
	"backend/seed-savers/utils"
	"fmt"
	"log"
	"strings"

	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.SeedStore
	usersStore   types.UserStore
	sessionStore *auth.AuthStore
}

func NewHandler(s types.SeedStore, us types.UserStore, sessionStore *auth.AuthStore) *Handler {
	return &Handler{s, us, sessionStore}
}

func (h *Handler) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/seeds", h.handleSeeds).Methods("GET")
	router.HandleFunc("/create-seed", auth.WithJWTAuth(h.handleCreateSeed, h.usersStore, h.sessionStore)).Methods("POST")
	router.HandleFunc("/update-seed", auth.WithJWTAuth(h.handleUpdateSeed, h.usersStore, h.sessionStore)).Methods("PUT")
}

func (h *Handler) handleSeeds(w http.ResponseWriter, r *http.Request) {
	seeds, err := h.store.GetSeeds()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, seeds)
}

func (h *Handler) handleCreateSeed(w http.ResponseWriter, r *http.Request) {

	//get the body of payload
	payload, err := utils.DecodePayload[types.CreateSeedPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//get userId
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	//check if seeds already exist in db
	seed, _ := h.store.GetSeedByVarieties(strings.ToLower(payload.Variety_name))
	if seed != nil {
		seed.Quantity = payload.Quantity
		err = h.usersStore.RegisterSeed(seed, userID)
		if err != nil {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("you have already registered this seed"))
		}
		utils.WriteJSON(w, http.StatusOK, nil)
		return
	}

	err = h.store.CreateSeed(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	seed, _ = h.store.GetSeedByVarieties(strings.ToLower(payload.Variety_name))

	seed.Quantity = payload.Quantity
	err = h.usersStore.RegisterSeed(seed, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Print(userID, payload)
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleUpdateSeed(w http.ResponseWriter, r *http.Request) {

	//get the body of payload
	payload, err := utils.DecodePayload[types.CreateSeedPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//get userId
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	//check if seeds already exist in db
	seed, _ := h.store.GetSeedByVarieties(strings.ToLower(payload.Variety_name))
	if seed == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("seed not exist, create one before udate it"))
		return
	}

	seed.Quantity = payload.Quantity

	h.usersStore.ModifySeedQuantity(seed, userID)
	utils.WriteJSON(w, http.StatusOK, nil)
}
