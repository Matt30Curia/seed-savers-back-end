package user

import (
	"backend/seed-savers/config"
	"backend/seed-savers/services/auth"
	"backend/seed-savers/services/email"
	"backend/seed-savers/types"
	"backend/seed-savers/utils"

	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	store        types.UserStore
	sessionStore *auth.AuthStore
}

func NewHandler(store types.UserStore, session *auth.AuthStore) *Handler {
	return &Handler{store: store, sessionStore: session}
}

func (h *Handler) RegisterRouter(router *mux.Router) {

	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/register/adress", auth.WithJWTAuth(h.handlePOSTAdress, h.store, h.sessionStore)).Methods("POST")
	router.HandleFunc("/register/adress", auth.WithJWTAuth(h.handlePUTAdress, h.store, h.sessionStore)).Methods("PUT")
	router.HandleFunc("/user/delete", auth.WithJWTAuth(h.handleDeleteUser, h.store, h.sessionStore)).Methods("DELETE")
	router.HandleFunc("/user/reset", h.handleResetSendEmail).Methods(http.MethodPost)
	router.HandleFunc("/user/reset/{encripted:.*}", h.handleResetPassword).Methods(http.MethodPost)

	//login whith google
	router.HandleFunc("/auth/{provider}", h.HandleAuthProvider).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", h.HandleAuthCallbackFunction).Methods("GET")
	router.HandleFunc("/auth/logout/{email}", h.HandleAuthLogOut).Methods("GET")

}

func (h *Handler) handleResetSendEmail(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.DecodePayload[types.UserRecoveryPassword](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	token, _ := utils.GetAESEncrypted(fmt.Sprintf("%v&%v", user.Email, time.Now().Format(layout)))

	subject := "Subject: Reset Password from Seed Savers\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><h1>clicca il link per il recupero</h1><a href='http://localhost:3000/user/reset/%v'>click this lin </a></body></html>", token)
	msg := []byte(subject + mime + body)

	if err := email.SendMail(user.Email, msg); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, nil)
}

func (h *Handler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.DecodePayload[types.UserResetPassword](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	encryptedEmailAndTime := mux.Vars(r)["encripted"]
	emailAndTime, err := utils.GetAESDecrypted(encryptedEmailAndTime)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	parts := strings.Split(string(emailAndTime), "&")
	mail := parts[0]
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	tIme, err := time.Parse(layout, parts[1])

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if uint8(time.Until(tIme).Abs().Hours()) >= uint8(config.Envs.TokenExpirationInHour) {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("token expired. Retry"))
		return
	}

	user, err := h.store.GetUserByEmail(string(mail))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	user.Password = hashPassword

	err = h.store.ModifyUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)

}

func (h *Handler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {

	id, err := auth.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	err = h.store.DeleteUserByID(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handlePOSTAdress(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.DecodePayload[types.AdressPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	err = h.store.CreateAdress(&types.Adress{
		ID:               id,
		Country:          payload.State,
		Street:           payload.Street,
		City:             payload.City,
		Cap:              payload.Cap,
		Province:         payload.Province,
		Apartment_number: payload.Apartment_number,
		Number:           payload.Number,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handlePUTAdress(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.DecodePayload[types.AdressPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}
	err = h.store.ModifyAdress(&types.Adress{
		ID:               id,
		Country:            payload.State,
		Street:           payload.Street,
		City:             payload.City,
		Cap:              payload.Cap,
		Province:         payload.Province,
		Apartment_number: payload.Apartment_number,
		Number:           payload.Number,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.DecodePayload[types.UserLoginPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) || u.Password == "OAUTH" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, uint64(u.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	h.sessionStore.StoreUserSession(w, r, token)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.DecodePayload[types.UserRegisterPayload](w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//check if the user exist
	_, err = h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with this email: %s already exist", payload.Email))
		return
	}

	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//if it doesn't exist create a new user
	err = h.store.CreateUser(&types.User{
		Email:    payload.Email,
		Password: hashPassword,
		Name:     payload.Name,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {

	provider := mux.Vars(r)["provider"]
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//check if user exist in db
	_, err = h.store.GetUserByEmail(user.Email)
	if err != nil {
		//if not exist create user in db
		h.store.CreateUser(&types.User{Email: user.Email, Name: user.NickName, Password: "OAUTH"})
	}

	//create autorization token
	u, _ := h.store.GetUserByEmail(user.Email)
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, uint64(u.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	h.sessionStore.StoreUserSession(w, r, token)

	http.Redirect(w, r, "http://localhost:5173/profile", http.StatusFound)
}

func (h *Handler) HandleAuthLogOut(w http.ResponseWriter, r *http.Request) {
	h.sessionStore.RemoveUserSession(w, r)
	gothic.Logout(w, r)
}

func (h *Handler) HandleAuthProvider(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	provider := mux.Vars(r)["provider"]
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}
