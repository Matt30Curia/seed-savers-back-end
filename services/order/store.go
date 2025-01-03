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
func (s *Store) GetIncomingOrders(reciverUserID int) ([]types.Order, error) {

	query := `SELECT o.*, reciver.name, a.state, a.city, a.street, a.cap, 
			  a.province, a.number, a.apartment_number, s.img, s.variety_name,  od.quantity, s.seed_id
			  FROM orders o
 			  JOIN users reciver ON o.reciver_user_id = reciver.user_id
			  LEFT JOIN adress a ON reciver.user_id = a.id
			  JOIN order_detail od ON o.order_id = od.order_id
			  JOIN seed s ON od.seed_id = s.seed_id
			  WHERE o.reciver_user_id = ?;`

	// Query per ottenere tutti gli ordini ricevuti da un utente specifico
	rows, err := s.db.Query(query, reciverUserID)
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
func (s *Store) GetOrdersToBeSent(senderUserID int) ([]types.Order, error) {

	query := `SELECT o.*, sender.name, a.state, a.city, a.street, a.cap, 
			  a.province, a.number, a.apartment_number, s.img, s.variety_name,  od.quantity, s.seed_id
			  FROM orders o
	 		  JOIN users sender ON o.sender_user_id = sender.user_id
			  LEFT JOIN adress a ON sender.user_id = a.id
			  JOIN order_detail od ON o.order_id = od.order_id
			  JOIN seed s ON od.seed_id = s.seed_id
			  WHERE o.sender_user_id = ?;`
	// Query per ottenere tutti gli ordini inviati da un utente specifico
	rows, err := s.db.Query(query, senderUserID)
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

func (s *Store) DeleteOrder(ID int) error {
	// Inizio della transazione
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Rollback automatico se qualcosa va storto
	defer tx.Rollback()

	// Inseriamo l'ordine nella tabella orders
	_, err = tx.Exec("DELETE FROM order_detail where order_id = ?;",ID)
	if err != nil {
		return err
	}

	// Inseriamo i dettagli dell'ordine nella tabella order_detail
	_, err = tx.Exec("DELETE FROM orders where order_id = ? ;", ID)
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
	// Create a new Order object to fill with row data
	order := new(types.Order)

	// Temporary variables for scanning
	var (
		orderID       int
		senderUserID  int
		reciverUserID int
		orderDate     time.Time
		orderState    string
		senderName    string
		country       string
		city          string
		street        string
		cap           string
		province      string
		aptNumber     string
		number        uint16
		img           string
		varietyName   string
		quantity      int
		seedId        int
	)

	// Scan the row into our variables - matched with the column order provided
	err := rows.Scan(
		&orderID,       // order_id
		&senderUserID,  // sender_user_id
		&reciverUserID, // reciver_user_id
		&orderDate,     // order_date
		&orderState,    // state
		&senderName,    // sender_name
		&country,       // state
		&city,          // city
		&street,        // street
		&cap,           // cap
		&province,      // province
		&number,        // number
		&aptNumber,     // apartment_number
		&img,           // img
		&varietyName,   // variety_name
		&quantity,      // quantity
		&seedId,
	)

	if err != nil {
		return nil, fmt.Errorf("error scanning order row: %w", err)
	}

	// Populate Order fields
	
	order.ID = orderID
	order.State = orderState
	order.OrderDate = orderDate
	order.SenderID = senderUserID
	order.ReciverID = reciverUserID

	// Construct the Address
	order.ReciverAdress = types.Adress{
		Street:           street,
		City:             city,
		Cap:              cap,
		Province:         province,
		Number:           number,
		Apartment_number: aptNumber,
		Country:          country, // Now using the actual scanned value
	}

	// Construct the Seed
	order.Seed = types.Seed{
		ID:           seedId,
		Image:        img,
		Variety_name: varietyName,
		Quantity:     quantity,
	}

	return order, nil
}
