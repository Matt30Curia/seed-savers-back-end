package auth

import (
	"backend/seed-savers/types"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const SessionName = "authorization"

type SessionOptions struct {
	CookiesKey string
	MaxAge     int
	HttpOnly   bool // Should be true if the site is served over HTTP (development environment)
	Secure     bool // Should be true if the site is served over HTTPS (production environment)
}

type AuthStore struct {
	Store sessions.Store
}

// CreateAdress implements types.UserStore.
func (authStore *AuthStore) CreateAdress(adress *types.Adress) error {
	panic("unimplemented")
}

// CreateUser implements types.UserStore.
func (authStore *AuthStore) CreateUser(user *types.User) error {
	panic("unimplemented")
}

// DeleteUserByID implements types.UserStore.
func (authStore *AuthStore) DeleteUserByID(ID int) error {
	panic("unimplemented")
}

// GetCompleteUserByEmail implements types.UserStore.
func (authStore *AuthStore) GetCompleteUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}

// GetCompleteUserByID implements types.UserStore.
func (authStore *AuthStore) GetCompleteUserByID(ID int) (*types.User, error) {
	panic("unimplemented")
}

// GetUserByEmail implements types.UserStore.
func (authStore *AuthStore) GetUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}

// GetUserByID implements types.UserStore.
func (authStore *AuthStore) GetUserByID(ID int) (*types.User, error) {
	panic("unimplemented")
}

// ModifyAdress implements types.UserStore.
func (authStore *AuthStore) ModifyAdress(adress *types.Adress) error {
	panic("unimplemented")
}

// ModifySeedQuantity implements types.UserStore.
func (authStore *AuthStore) ModifySeedQuantity(seed *types.Seed, userID int) error {
	panic("unimplemented")
}

// ModifyUser implements types.UserStore.
func (authStore *AuthStore) ModifyUser(user *types.User) error {
	panic("unimplemented")
}

// RegisterSeed implements types.UserStore.
func (authStore *AuthStore) RegisterSeed(seed *types.Seed, userID int) error {
	panic("unimplemented")
}

func NewCookieStore(opts SessionOptions) *AuthStore {
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))

	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	return &AuthStore{Store: store}
}

func (authStore *AuthStore) StoreUserSession(w http.ResponseWriter, r *http.Request, userID string) error {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := authStore.Store.Get(r, SessionName)

	session.Values["user"] = userID

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func (authStore *AuthStore) GetSessionUserToken(r *http.Request) (string, error) {

	session, err := authStore.Store.Get(r, SessionName)
	if err != nil {
		return "", err
	}

	u := session.Values["user"]
	if u == nil {
		return "", fmt.Errorf("user is not authenticated! %v", u)
	}

	return u.(string), nil
}

func (authStore *AuthStore) RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	session, err := authStore.Store.Get(r, SessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = ""
	// delete the cookie immediately
	session.Options.MaxAge = -1

	session.Save(r, w)
}
