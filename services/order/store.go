package order

import (
	"backend/seed-savers/types"
	"database/sql"
	"fmt"
	"time"
)

// Store rappresenta una struttura che gestisce l'accesso al database per gli ordini
type Store struct {
	db *sql.DB
}

// NewStore crea e restituisce un nuovo oggetto Store con il database passato come parametro
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetOrdersById restituisce un ordine dato il suo ID
func (s *Store) GetOrdersById(ID int) (*types.Order, error) {
	// Eseguiamo la query per ottenere l'ordine tramite ID
	rows, err := s.db.Query("SELECT * FROM orders WHERE order_id=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Creiamo una nuova struttura ordine
	order := new(types.Order)
	
	// Iteriamo sulle righe della query
	for rows.Next() {
		order, err = ScanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
	}

	// Se l'ordine non è stato trovato, ritorniamo un errore
	if order.ID == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

// GetOrdersByReciver restituisce una lista di ordini ricevuti da un utente dato l'ID
func (s *Store) GetOrdersByReciver(reciverUserID int) ([]types.Order, error) {
	// Query per ottenere tutti gli ordini ricevuti da un utente specifico
	rows, err := s.db.Query("SELECT * FROM orders WHERE reciver_user_id=?", reciverUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0)

	// Iteriamo sulle righe restituite dalla query
	for rows.Next() {
		order, err := ScanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	// Se non ci sono ordini, ritorniamo un errore
	if len(orders) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return orders, nil
}

// GetOrdersBySender restituisce una lista di ordini inviati da un utente dato l'ID
func (s *Store) GetOrdersBySender(senderUserID int) ([]types.Order, error) {
	// Query per ottenere tutti gli ordini inviati da un utente specifico
	rows, err := s.db.Query("SELECT * FROM orders WHERE sender_user_id=?", senderUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0)

	// Iteriamo sulle righe restituite dalla query
	for rows.Next() {
		order, err := ScanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	// Se non ci sono ordini, ritorniamo un errore
	if len(orders) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return orders, nil
}

// MakeOrder crea un nuovo ordine e i dettagli associati (come quantità e seme) con una transazione
func (s *Store) MakeOrder(senderUserID, reciverUserID, seedId, quantity int) error {
	// Inizio della transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Rollback automatico se qualcosa va storto
	defer tx.Rollback()

	// Inseriamo l'ordine nella tabella orders
	res, err := tx.Exec("INSERT INTO orders (sender_user_id, reciver_user_id) VALUES (?, ?)", senderUserID, reciverUserID)
	if err != nil {
		return err
	}

	// Otteniamo l'ID dell'ordine appena creato
	orderID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Inseriamo i dettagli dell'ordine nella tabella order_detail
	_, err = tx.Exec("INSERT INTO order_detail (order_id, seed_id, quantity) VALUES (?, ?, ?)", orderID, seedId, quantity)
	if err != nil {
		return err
	}

	// Confermiamo la transazione
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// ModifyOrder modifica un ordine esistente, come lo stato e la quantità di semi
func (s *Store) ModifyOrder(order *types.Order) error {
	// Inizio della transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Rollback automatico se qualcosa va storto
	defer tx.Rollback()

	// Modifica lo stato dell'ordine nella tabella orders
	_, err = tx.Exec("UPDATE orders SET state = ? WHERE order_id = ?", order.State, order.ID)
	if err != nil {
		return err
	}

	// Modifica la quantità dei semi nei dettagli dell'ordine
	_, err = tx.Exec("UPDATE order_detail SET quantity = ? WHERE order_id = ?", order.Seed.Quantity, order.ID)
	if err != nil {
		return err
	}

	// Confermiamo la transazione
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// ScanRowIntoOrder esegue il binding dei dati di una riga su un oggetto Order
func ScanRowIntoOrder(rows *sql.Rows) (*types.Order, error) {
	// Creiamo un nuovo oggetto Order da riempire con i dati della riga
	order := new(types.Order)

	// Variabili temporanee per i campi letti dalla query
	var senderUserID, reciverUserID, seedId int
	var state string
	var orderDate time.Time

	// Effettuiamo la scansione della riga
	err := rows.Scan(
		&order.ID,           // order_id
		&senderUserID,       // sender_user_id
		&reciverUserID,      // reciver_user_id
		&orderDate,          // order_date
		&state,              // state (è una stringa che rappresenta un ENUM)
	)
	if err != nil {
		return nil, err
	}

	// Popoliamo i campi letti dalla query
	order.OrderDate = orderDate
	order.State = state

	// Popoliamo i campi Reciver e Sender con i dati degli utenti
	order.Sender = types.User{ID: senderUserID}
	order.Reciver = types.User{ID: reciverUserID}

	// Popoliamo il seme associato all'ordine
	order.Seed = types.Seed{ID: seedId}

	return order, nil
}
