package seed

import (
	"backend/seed-savers/types"
	"database/sql"
	"fmt"
	"strings"
)

// Store rappresenta una struttura che gestisce l'accesso al database per i semi
type Store struct {
	db *sql.DB
}

// NewStore crea e restituisce un nuovo oggetto Store con il database passato come parametro
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetSeeds restituisce una lista di tutti i semi nel database
func (s *Store) GetSeeds() ([]types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seeds := make([]types.Seed, 0, 3)

	// Itera sulle righe restituite dalla query e popola il slice di semi
	for rows.Next() {
		seed, err := ScanRowIntoSeed(rows)
		if err != nil {
			return nil, err
		}
		seeds = append(seeds, *seed)
	}

	if len(seeds) == 0 {
		return nil, fmt.Errorf("seeds not found")
	}
	return seeds, nil
}

// GetSeedByID restituisce un seme specifico dato il suo ID
func (s *Store) GetSeedByID(id int) (*types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seed := new(types.Seed)

	for rows.Next() {
		seed, err = ScanRowIntoSeed(rows)
		if err != nil {
			return nil, err
		}
	}

	if seed.ID == 0 {
		return nil, fmt.Errorf("seed not found")
	}

	return seed, nil
}

// GetTotalSeedQuantity restituisce la quantità totale di un seme, sommando le quantità per tutti gli utenti
func (s *Store) GetTotalSeedQuantity(id int) (int, error) {
	rows, err := s.db.Query("SELECT SUM(quantity) AS total_quantity FROM users_seed WHERE seed_id = ?", id)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var quantity int

	for rows.Next() {
		err := rows.Scan(&quantity)
		if err != nil {
			return -1, err
		}
	}
	return quantity, nil
}

// GetSeedOwnersByID restituisce una mappa con i nomi degli utenti e le quantità di semi che possiedono
func (s *Store) GetSeedOwnersByID(id int) (map[string]int, error) {
	rows, err := s.db.Query("SELECT us.seed_id, u.name, us.quantity FROM users_seed us INNER JOIN users u ON us.user_id = u.user_id WHERE us.seed_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	owners := make(map[string]int)

	var username string
	var quantity int

	// Popola la mappa con i dati di ogni proprietario di semi
	for rows.Next() {
		err := rows.Scan(new(int), &username, &quantity)
		if err != nil {
			return nil, err
		}
		owners[username] = quantity
	}
	return owners, nil
}

// GetSeedByVarieties restituisce un seme che corrisponde al nome della varietà
func (s *Store) GetSeedByVarieties(varieties string) (*types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed.variety_name LIKE ?", "%"+varieties+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seed *types.Seed

	for rows.Next() {
		seed, err = ScanRowIntoSeed(rows)
		if err != nil {
			return nil, err
		}
	}

	return seed, nil
}

// GetSeedsByVegetable restituisce una lista di semi che corrispondono a un determinato tipo di ortaggio
func (s *Store) GetSeedsByVegetable(vegetable string) ([]types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed.vegetable LIKE ?", "%"+vegetable+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seeds := make([]types.Seed, 0)

	// Aggiunge i semi alla lista
	for rows.Next() {
		seed, err := ScanRowIntoSeed(rows)
		if err != nil {
			return nil, err
		}
		seeds = append(seeds, *seed)
	}

	if len(seeds) == 0 {
		return nil, fmt.Errorf("seeds not found")
	}
	return seeds, nil
}

// CreateSeed crea un nuovo seme nel database
func (s *Store) CreateSeed(seedPayload *types.CreateSeedPayload) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Assicura che il rollback venga eseguito in caso di errore

	// Inserisce il seme nel database
	_, err = tx.Exec("INSERT INTO seed (description, variety_name, vegetable, img) VALUES (?, ?, ?, ?)", seedPayload.Description, strings.ToLower(seedPayload.Variety_name), seedPayload.Vegetable, seedPayload.Image)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

// UpdateSeed aggiorna le informazioni di un seme esistente
func (s *Store) UpdateSeed(seed *types.Seed) error {
	// Inizia una transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Assicura che il rollback venga eseguito in caso di errore

	// Modifica i dettagli del seme
	_, err = tx.Exec("UPDATE seed SET description=?, variety_name=?, vegetable=?, img=? WHERE seed_id=?", seed.Description, seed.Variety_name, seed.Vegetable, seed.Image, seed.ID)
	if err != nil {
		return err
	}

	// Conferma la transazione
	return tx.Commit()
}

func (s *Store) UserSeedQuantity(id, seedId int) int{
	rows, err := s.db.Query("SELECT quantity FROM users_seed WHERE user_id = ? AND seed_id = ?", id, seedId)
	if err != nil {
		return -1
	}
	defer rows.Close()

	var quantity int
	rows.Next()
	err = rows.Scan(&quantity)
	if err != nil || quantity <= 0{
		return -1
	}
	
	return quantity
}


// ScanRowIntoSeed esegue il binding dei dati di una riga su un oggetto Seed
func ScanRowIntoSeed(rows *sql.Rows) (*types.Seed, error) {
	seed := new(types.Seed)
	var img sql.NullString
	err := rows.Scan(
		&seed.ID,
		&seed.Description,
		&img,
		&seed.Variety_name,
		&seed.Vegetable,
	)
	if err != nil {
		return nil, err
	}

	// Gestione dell'immagine, che potrebbe essere null
	if img.Valid {
		seed.Image = img.String
	}

	return seed, nil
}