package seed

import (
	"backend/seed-savers/services/auth"
	"backend/seed-savers/types"

	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func TestUserServiceHandlers(t *testing.T) {

	mockStore := &mockUserStore{}
	autMockStore := &auth.AuthStore{Store: sessions.NewCookieStore([]byte{5})}
	sessionsMock := auth.NewCookieStore(auth.SessionOptions{})
	handler := NewHandler(mockStore, autMockStore, sessionsMock)

	t.Run("should return seeds list successfully", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/seeds", nil)
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/seeds", handler.handleSeeds)
		router.ServeHTTP(rr, req)

		// Checking the status code
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, rr.Code)
		}

	
		// Checking if response body contains the expected mock data
		expected := types.Seed{
			ID:           1,
			Variety_name: "Tomato",
			Description:  "Cherry tomatoes",
			Vegetable:    "Tomato",
			Image:        "image.com",
			Quantity:     100,
		}
		buf := new(bytes.Buffer)

		// Converte la struct in un array di byte in formato binario
		err = binary.Write(buf, binary.LittleEndian, expected)

		res, _ := json.Marshal(rr.Body.String())
		if reflect.DeepEqual(res, err){
			t.Errorf("expected body %v but got %s", expected, rr.Body.String())
		}
	})

	t.Run("should fail if the seed creation is invalid", func(t *testing.T) {
		payload := types.CreateSeedPayload{
			Description:  "",
			Variety_name: "",
			Vegetable:    "",
			Image:        "",
			Quantity:     0,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/create-seed", bytes.NewBuffer(marshalled))
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/create-seed", handler.handleCreateSeed)
		router.ServeHTTP(rr, req)

		// In case of invalid data, expect a bad request error (400)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserStore struct{}

// CreateSeed implements types.SeedStore.
func (m *mockUserStore) CreateSeed(*types.CreateSeedPayload) error {
	panic("unimplemented")
}

// GetSeedByID implements types.SeedStore.
func (m *mockUserStore) GetSeedByID(id int) (*types.Seed, error) {
	panic("unimplemented")
}

// GetSeedByVarieties implements types.SeedStore.
func (m *mockUserStore) GetSeedByVarieties(varieties string) (*types.Seed, error) {
	panic("unimplemented")
}

// GetSeedOwnersByID implements types.SeedStore.
func (m *mockUserStore) GetSeedOwnersByID(id int) (map[int]string, error) {
	panic("unimplemented")
}

// GetSeeds implements types.SeedStore.
func (m *mockUserStore) GetSeeds() ([]types.Seed, error) {
	return []types.Seed{
		{
			ID:           1,
			Variety_name: "Tomato",
			Description:  "Cherry tomatoes",
			Vegetable:    "Tomato",
			Image:        "image.com",
			Quantity:     100,
		},
	}, nil
}

// GetSeedsByVegetable implements types.SeedStore.
func (m *mockUserStore) GetSeedsByVegetable(vegetable string) ([]types.Seed, error) {
	panic("unimplemented")
}

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
