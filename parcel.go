package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

// NewParcelStore инициализирует новый ParcelStore
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет новую посылку в базу данных
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		`INSERT INTO parcel (client, status, address, created_at) 
		VALUES (?, ?, ?, ?)`,
		p.Client,
		p.Status,
		p.Address,
		p.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add parcel: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return int(id), nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Чтение строк из таблицы по client
	rows, err := s.db.Query(`SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get parcels by client: %w", err)
	}
	defer rows.Close()

	// Заполнение среза Parcel данными из таблицы
	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan parcel: %w", err)
		}
		parcels = append(parcels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return parcels, nil
}

// Delete удаляет посылку, если она в статусе "registered"
func (s ParcelStore) Delete(number int) error {
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("can only delete registered parcel")
	}
	_, err = s.db.Exec(`DELETE FROM parcel WHERE number = ?`, number)
	return err
}

// SetAddress обновляет адрес посылки, если она в статусе "registered"
func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("failed to get parcel for update: %w", err)
	}
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("can only change address of registered parcel")
	}
	_, err = s.db.Exec(`UPDATE parcel SET address = ? WHERE number = ?`, address, number)
	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}
	return nil
}

// SetStatus обновляет статус посылки
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(`UPDATE parcel SET status = ? WHERE number = ?`, status, number)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

// Get возвращает посылку по её номеру
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow(`SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`, number)
	var p Parcel
	if err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return Parcel{}, fmt.Errorf("parcel not found: %w", err)
		}
		return Parcel{}, fmt.Errorf("failed to scan parcel: %w", err)
	}
	return p, nil
}
