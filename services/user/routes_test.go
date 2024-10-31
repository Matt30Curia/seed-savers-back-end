package user

import (
	"backend/seed-savers/services/auth"
	"backend/seed-savers/types"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func TestUserServiceHandlers(t *testing.T) {

	mockStore := &mockUserStore{}
	autMockStore := &auth.AuthStore{Store: sessions.NewCookieStore([]byte{5})}
	handler := NewHandler(mockStore, autMockStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.UserRegisterPayload{
			Name:     "matteo",
			Email:    "invalid",
			Password: "asd",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expeted status code %d but got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail if the user creation is invalid", func(t *testing.T) {
		payload := types.UserRegisterPayload{
			Name:     "matteo",
			Email:    "valid@mail.com",
			Password: "asd",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expeted status code %d but got %d", http.StatusCreated, rr.Code)
		}
	})

}

type mockUserStore struct{}

// GetCompleteUserByEmail implements types.UserStore.
func (m *mockUserStore) GetCompleteUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}

// GetCompleteUserByID implements types.UserStore.
func (m *mockUserStore) GetCompleteUserByID(ID int) (*types.User, error) {
	panic("unimplemented")
}

// ModifySeedQuantity implements types.UserStore.
func (m *mockUserStore) ModifySeedQuantity(seed *types.Seed, userID int) error {
	panic("unimplemented")
}

func (m *mockUserStore) UpdateUser(user *types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, fmt.Errorf("user not found")
}

func (m *mockUserStore) CreateUser(user *types.User) error {
	return nil
}

func (m *mockUserStore) ModifyUser(user *types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) RegisterSeed(seed *types.Seed, userID int) error {

	return nil
}

func (m *mockUserStore) CreateAdress(adress *types.Adress) error {
	return nil
}

func (m *mockUserStore) ModifyAdress(adress *types.Adress) error {
	return nil
}

func (m *mockUserStore) DeleteUserByID(ID int) error {
	return nil
}
