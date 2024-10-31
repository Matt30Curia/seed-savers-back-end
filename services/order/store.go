package order

import (
	"backend/seed-savers/types"
	"database/sql"
	"fmt"
	//"log"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}


func (s *Store) GetOrdersById(ID int) (*types.Order, error) {
	rows, err := s.db.Query("SELECT * FROM seed WHERE seed_id=?", ID)
	if err != nil {
		return nil, err
	}

	order := new(types.Order)
	
	for rows.Next() {
		order, err = ScanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
	}

	if order.ID == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

func (s *Store) GetOrdersByReciver(reciverUserID int) ([]types.Order, error){
	rows, err := s.db.Query("SELECT * FROM orders WHERE reciver_user_id=?", reciverUserID)
	if err != nil {
		return nil, err
	}

	orders := make([]types.Order, 10)
	var order *types.Order
	for rows.Next() {
		order, err = ScanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	
	if order.ID == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return orders, nil
}

func (s *Store) GetOrdersBySender(senderUserID int) ([]types.Order, error){
	rows, err := s.db.Query("SELECT * FROM orders WHERE sender_user_id=?", senderUserID)
	if err != nil {
		return nil, err
	}

	orders := make([]types.Order, 10)
	var order *types.Order
	for rows.Next() {
		order, err = ScanRowIntoOrder(rows)
		if err != nil { 
			return nil, err
		}
		orders = append(orders, *order)
	}

	if order.ID == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return orders, nil
}

//TO-DO inser Quantity in order detail
func (s *Store) MakeOrder(senderUserID, reciverUserID, seedId, quantity int) error {
    // Inizio della transazione
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }

    // Inserisci l'ordine nella tabella `orders`
    res, err := tx.Exec("INSERT INTO orders (sender_user_id, reciver_user_id) VALUES (?, ?)", senderUserID, reciverUserID)
    if err != nil {
        tx.Rollback() // Rollback in caso di errore
        return err
    }

    // Ottieni l'order_id appena creato
    orderID, err := res.LastInsertId()
    if err != nil {
        tx.Rollback() // Rollback in caso di errore
        return err
    }
	fmt.Println(orderID)
    // Inserisci i dettagli dell'ordine nella tabella `order_detail`
    _, err = tx.Exec("INSERT INTO order_detail (order_id, seed_id, quantity) VALUES (?, ?, ?)", orderID, seedId, quantity)
    if err != nil {
        tx.Rollback() // Rollback in caso di errore
        return err
    }

    // Conferma la transazione
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}



func (s *Store) ModifyOrder(order *types.Order) error {
    // Inizio della transazione
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    
    // Inserisci l'ordine nella tabella `orders`
    _, err = tx.Exec(
		"UPDATE orders SET sender_user_id = ?, reciver_user_id = ?, state = ? WHERE order_id = ?;",
		order.Sender.ID, order.Reciver.ID, order.State, order.ID)	
	if err != nil {
        tx.Rollback() // Rollback in caso di errore
        return err
    }

    // Inserisci i dettagli dell'ordine nella tabella `order_detail`
    _, err = tx.Exec("INSERT INTO order_detail (order_id, seed_id, quantity) VALUES (?, ?, ?)", order.ID, order.Seed.ID, order.Seed.Quantity)
    if err != nil {
        tx.Rollback() // Rollback in caso di errore
        return err
    }

    // Conferma la transazione
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}


func ScanRowIntoOrder(rows *sql.Rows) (*types.Order, error) {
    // Creiamo un nuovo oggetto Order da riempire con i dati della riga
    order := new(types.Order)

    // Dichiarazione di variabili temporanee per i valori che verranno letti
    var senderUserID, reciverUserID, seedId int
    var state string
	var orderDate time.Time
    // Effettuiamo la scansione della riga principale
    err := rows.Scan(
        &order.ID,           // order_id
        &senderUserID,       // sender_user_id
        &reciverUserID,      // reciver_user_id
        &orderDate,       // order_date
        &state,              // state (è un ENUM quindi una stringa in Go)

    )
    if err != nil {
        return nil, err
    }

    // Popoliamo i campi letti
    order.OrderDate = orderDate
    order.State = state

    // Popoliamo i campi Reciver e Sender con i dati degli utenti
    order.Sender = types.User{ID : senderUserID}
    order.Reciver = types.User{ID : reciverUserID}

    
    order.Seed = types.Seed{ID: seedId}

    return order, nil
}

