package seed

import (
	"backend/seed-savers/types"
	"database/sql"
	"fmt"

	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetSeeds() ([]types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed")

	if err != nil {
		return nil, err
	}

	seeds := make([]types.Seed, 0, 3)

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

func (s *Store) GetSeedByID(id int) (*types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed_id=?", id)

	if err != nil {
		return nil, err
	}

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

func (s *Store) GetTotalSeedQuantity(id int) (int, error){
	rows, err :=s.db.Query("SELECT SUM(quantity) AS total_quantity FROM users_seed WHERE seed_id = ?;", id)
	
	if err != nil {
		return -1, err
	}

	var quantity int

	for rows.Next() {
		err := rows.Scan(
			&quantity,
		)

		if err != nil {
			return -1, err
		}
	
	}
	return quantity, nil
}

func (s *Store) GetSeedOwnersByID(id int) (map[int]string, error) {	
	rows, err := s.db.Query("SELECT us.user_id, us.seed_id, u.name, us.quantity FROM users_seed us INNER JOIN users u ON us.user_id = u.user_id where us.seed_id = ?;", id)

	if err != nil {
		return nil, err
	}

	owners := make(map[int]string)

	ID := -1
	user_id := ""

	for rows.Next() {
		err := rows.Scan(
			&ID,
			&user_id,
		)

		if err != nil {
			return nil, err
		}
		owners[ID] = user_id
	}
	return owners, nil
}

func (s *Store) GetSeedByVarieties(varieties string) (*types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed.variety_name LIKE ?", "%"+varieties+"%")

	if err != nil {
		return nil, err
	}

	var seed *types.Seed

	for rows.Next() {
		seed, err = ScanRowIntoSeed(rows)
		if err != nil {
			return nil, err
		}
	}

	return seed, nil
}

func (s *Store) GetSeedsByVegetable(vegetable string) ([]types.Seed, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed.variety_name LIKE '%%%s%%'", vegetable)

	if err != nil {
		return nil, err
	}

	seeds := make([]types.Seed, 1)

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

func (s *Store) CreateSeed(seedPayload *types.CreateSeedPayload) error {
	_, err := s.db.Exec("INSERT INTO seed (description, variety_name, vegetable, img) VALUES (?, ?, ?, ?)", seedPayload.Description, strings.ToLower(seedPayload.Variety_name), seedPayload.Vegetable, seedPayload.Image)
	
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateSeed(seed *types.Seed) error {
	_, err := s.db.Exec("UPDATE seed SET description=?, SET variety_name=?, SET vegetable=?, SET img=? WHERE seed.seed_id = ?", seed.Description, seed.Variety_name, seed.Vegetable, seed.Image, seed.ID)
	
	if err != nil {
		return err
	}

	return nil
}

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

	if img.Valid {
		seed.Image = img.String
	}

	return seed, nil
}

/*
SELECT
    us.user_id,
    us.seed_id,
    u.name,
    us.quantity
FROM
    users_seed us
INNER JOIN
    users u ON us.user_id = u.user_id
where us.seed_id = 4;
*/

/*SELECT
    s.seed_id,
    s.img,
    s.variety_name,
    s.description,
    s.vegetable,
    us.user_id,
    u.name,
    us.quantity
FROM
    seed s

INNER JOIN
    users_seed us ON s.seed_id = us.seed_id
INNER JOIN
    users u ON us.user_id = u.user_id
where s.seed_id = !!!!!il tuo id varieta;*/
