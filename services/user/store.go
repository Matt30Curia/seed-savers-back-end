package user

import (
	"backend/seed-savers/types"
	"database/sql"
	"fmt"
)

// Store rappresenta una struttura per l'accesso al database
type Store struct {
	db *sql.DB
}

// NewStore crea e restituisce un nuovo oggetto Store con il database passato come parametro
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail cerca un utente nel database usando l'email e restituisce l'utente trovato
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(types.User)
	for rows.Next() {
		u, err = ScanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

// RegisterSeed registra un nuovo seme associato a un utente
func (s *Store) RegisterSeed(seed *types.Seed, userID int) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Assicurati di eseguire il rollback in caso di errore

	// Inserisci il seme associato all'utente
	_, err = tx.Exec("INSERT INTO users_seed (seed_id, user_id, quantity) VALUES (?, ?, ?)", seed.ID, userID, seed.Quantity)
	if err != nil {
		return err
	}

	// Se tutto è andato bene, conferma la transazione
	return tx.Commit()
}

// GetUserByID cerca un utente nel database usando l'ID e restituisce l'utente trovato
func (s *Store) GetUserByID(ID int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE user_id=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(types.User)
	for rows.Next() {
		u, err = ScanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

// CreateUser crea un nuovo utente nel database
func (s *Store) CreateUser(user *types.User) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Esegui rollback in caso di errore

	// Inserisce l'utente
	_, err = tx.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// ModifyUser modifica un utente esistente nel database
func (s *Store) ModifyUser(user *types.User) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Esegui rollback in caso di errore

	// Modifica l'utente
	_, err = tx.Exec("UPDATE users SET name = ?, email = ?, credits = ? ,password = ? WHERE user_id = ?", user.Name, user.Email, user.Credits, user.Password, user.ID)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// DeleteUserByID elimina un utente dal database usando l'ID
func (s *Store) DeleteUserByID(ID int) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Esegui rollback in caso di errore

	// Elimina l'utente
	_, err = tx.Exec("DELETE FROM users WHERE user_id = ?", ID)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// CreateAdress crea un nuovo indirizzo per un utente
func (s *Store) CreateAdress(adress *types.Adress) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Esegui rollback in caso di errore

	// Inserisci l'indirizzo
	_, err = tx.Exec("INSERT INTO adress (id, state, city, street, cap, province, number, apartment_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", adress.ID, adress.State, adress.City, adress.Street, adress.Cap, adress.Province, adress.Number, adress.Apartment_number)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// ModifyAdress modifica un indirizzo esistente nel database
func (s *Store) ModifyAdress(adress *types.Adress) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Esegui rollback in caso di errore

	// Modifica l'indirizzo
	_, err = tx.Exec("UPDATE adress SET state = ?, city = ?, street = ?, cap = ?, province = ?, number = ?, apartment_number = ? WHERE id = ?", adress.State, adress.City, adress.Street, adress.Cap, adress.Province, adress.Number, adress.Apartment_number, adress.ID)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// GetCompleteUserByEmail restituisce un utente con tutti i dettagli (indirizzo e semi) usando l'email
func (s *Store) GetCompleteUserByEmail(email string) (*types.User, error) {
	query := `
		SELECT u.user_id, u.NAME, u.email, u.credits, 
		       a.state, a.city, a.street, a.cap, a.province, a.number, a.apartment_number, 
		       s.seed_id, s.variety_name, s.description, s.vegetable, s.img, us.quantity 
		FROM users u
		LEFT JOIN adress a ON u.user_id = a.id
		LEFT JOIN users_seed us ON u.user_id = us.user_id
		LEFT JOIN seed s ON us.seed_id = s.seed_id
		WHERE u.email = ?;
	`
	rows, err := s.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := ScanRowsIntoCompleteUser(rows)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetCompleteUserByID restituisce un utente con tutti i dettagli (indirizzo e semi) usando l'ID
func (s *Store) GetCompleteUserByID(ID int) (*types.User, error) {
	query := `
		SELECT u.user_id, u.NAME, u.email, u.credits, 
		       a.state, a.city, a.street, a.cap, a.province, a.number, a.apartment_number, 
		       s.seed_id, s.variety_name, s.description, s.vegetable, s.img, us.quantity 
		FROM users u
		LEFT JOIN adress a ON u.user_id = a.id
		LEFT JOIN users_seed us ON u.user_id = us.user_id
		LEFT JOIN seed s ON us.seed_id = s.seed_id
		WHERE u.user_id = ?;
	`
	rows, err := s.db.Query(query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := ScanRowsIntoCompleteUser(rows)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ModifySeedQuantity aggiorna la quantità di semi per un utente
func (s *Store) ModifySeedQuantity(seed *types.Seed, userID int) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Esegui rollback in caso di errore

	// Aggiorna la quantità del seme
	_, err = tx.Exec("UPDATE users_seed SET quantity = ? WHERE user_id = ? AND seed_id = ?", seed.Quantity, userID, seed.ID)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// ScanRowIntoUser esegue il binding dei dati di una riga su un oggetto User
func ScanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Credits,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func ScanRowsIntoCompleteUser(rows *sql.Rows) (*types.User, error) {
	// Variabile per memorizzare l'utente
	var user *types.User
	var firstRow = true

	// Itera su tutte le righe
	for rows.Next() {
		// Creiamo temporaneamente oggetti Seed e Adress per ogni riga
		var seed types.Seed
		var address types.Adress

		// Estrazione dati dalla riga
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Credits,
			&address.ID,
			&address.State,
			&address.City,
			&address.Street,
			&address.Cap,
			&address.Province,
			&address.Apartment_number,
			&address.Number,
			&seed.Variety_name,
			&seed.Description,
			&seed.Vegetable,
			&seed.Image,
			&seed.Quantity,
		)

		if err != nil {
			return nil, err
		}

		// Se è la prima riga, inizializza l'utente e assegna l'indirizzo
		if firstRow {
			user = new(types.User)
			user.Adress = address
			firstRow = false
		}

		// Aggiunge il seme alla lista dei semi dell'utente
		user.Seeds = append(user.Seeds, seed)
	}

	if user == nil {
		return nil, fmt.Errorf("no user found")
	}

	return user, nil
}
