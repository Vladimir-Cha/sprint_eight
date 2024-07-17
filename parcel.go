package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

// NewParcelStore создаёт новое хранилище посылок
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// AddParcel добавляет новую посылку в базу данных
func (store ParcelStore) AddParcel(parcel Parcel) (int, error) {
	var id int
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?) RETURNING number`
	err := store.db.QueryRow(query, parcel.Client, parcel.Status, parcel.Address, parcel.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetParcel получает посылку из базы данных по идентификатору
func (store ParcelStore) GetParcel(id int) (Parcel, error) {
	var parcel Parcel
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`
	err := store.db.QueryRow(query, id).Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return parcel, nil
}

// DeleteParcel удаляет посылку из базы данных по идентификатору
func (store ParcelStore) DeleteParcel(id int) error {
	query := `DELETE FROM parcel WHERE number = ?`
	_, err := store.db.Exec(query, id)
	return err
}

// SetAddress обновляет адрес посылки в базе данных
func (store ParcelStore) SetAddress(id int, address string) error {
	query := `UPDATE parcel SET address = ? WHERE number = ?`
	_, err := store.db.Exec(query, address, id)
	return err
}

// SetStatus обновляет статус посылки в базе данных
func (store ParcelStore) SetStatus(id int, status string) error {
	query := `UPDATE parcel SET status = ? WHERE number = ?`
	_, err := store.db.Exec(query, status, id)
	return err
}

// GetParcelsByClient получает список посылок по идентификатору клиента
func (store ParcelStore) GetParcelsByClient(client int) ([]Parcel, error) {
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`
	rows, err := store.db.Query(query, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var parcel Parcel
		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, parcel)
	}
	return parcels, nil
}
