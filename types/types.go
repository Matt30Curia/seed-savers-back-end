package types

import "time"

type UserResetPassword struct {
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type UserRecoveryPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type UserRegisterPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdressPayload struct {
	State            string `json:"state" validate:"required"`
	City             string `json:"city" validate:"required"`
	Street           string `json:"street" validate:"required"`
	Cap              string `json:"cap" validate:"required"`
	Province         string `json:"province" validate:"required"`
	Apartment_number string `json:"apartment_number"`
	Number           uint16 `json:"number" validate:"required"`
}

type CreateSeedPayload struct {
	Description  string `json:"description" validate:"required"`
	Variety_name string `json:"variety_name" validate:"required"`
	Vegetable    string `json:"vegetable" validate:"required"`
	Image        string `json:"image" validate:"required"`
	Quantity     int    `json:"quantity" validate:"required"`
}

type UpdateOrderPayload struct {
	OrderId      int    `json:"orderId" validate:"required"`
	SeedQuantity int    `json:"seedQuantity"`
	State        string `json:"state"`
}

type OrderPayload struct {
	SenderID     int `json:"sender" validate:"required"`
	SeedID       int `json:"seedId" validate:"required"`
	SeedQuantity int `json:"seedQuantity"`
}

type User struct {
	Name     string `json:"firstName"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Adress   Adress `json:"adress"`
	Seeds    []Seed `json:"seeds"`
	Credits  int8   `json:"credits"`
	ID       int    `json:"id"`
}

type Adress struct {
	ID               int    `json:"id"`
	Country          string `json:"state"`
	City             string `json:"city"`
	Street           string `json:"street"`
	Cap              string `json:"cap"`
	Province         string `json:"province"`
	Apartment_number string `json:"apartment_number"`
	Number           uint16 `json:"Number"`
}

type Seed struct {
	Description  string `json:"description"`
	Variety_name string `json:"variety_name"`
	Vegetable    string `json:"vegetable"`
	Image        string `json:"image"`
	Quantity     int    `json:"quantity"`
	ID           int    `json:"id"`
}

type Order struct {
	ID            int       `json:"order_id"`
	State         string    `json:"state"`
	OrderDate     time.Time `json:"order-date"`
	ReciverID     int       `json:"reciverID"`
	ReciverName   string    `json:"reciverName"`
	ReciverAdress Adress    `json:"adress"`
	SenderID      int       `json:"senderID"`
	SenderName    string    `json:"senderName"`
	Seed          Seed      `json:"seed"`
}

type OrderStore interface {
	GetOrdersById(ID int) (*Order, error)
	GetIncomingOrders(reciverUserID int) ([]Order, error)
	GetOrdersToBeSent(senderUserID int) ([]Order, error)
	MakeOrder(reUserID, reciverUserID, seedId, quantity int) error
	ModifyOrder(order *Order) error
	DeleteOrder(ID int) error
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(ID int) (*User, error)
	DeleteUserByID(ID int) error
	CreateUser(user *User) error
	ModifyUser(user *User) error
	GetCompleteUserByEmail(email string) (*User, error)
	GetCompleteUserByID(ID int) (*User, error)

	CreateAdress(adress *Adress) error
	ModifyAdress(adress *Adress) error

	RegisterSeed(seed *Seed, userID int) error
	ModifySeedQuantity(seed *Seed, userID int) error
}

type SeedStore interface {
	GetSeeds() ([]Seed, error)
	GetSeedByID(id int) (*Seed, error)
	GetSeedByVarieties(varieties string) (*Seed, error)
	GetSeedsByVegetable(vegetable string) ([]Seed, error)
	CreateSeed(*CreateSeedPayload) error
	GetSeedOwnersByID(id int) (map[string]int, error)
	UserSeedQuantity(id, seedId int) int
}
