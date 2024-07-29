package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// Добавляет новую посылку в базу данных
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)`
	res, err := s.db.Exec(query, sql.Named("client", p.Client), sql.Named("status", p.Status), sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// Получает посылку по её номеру
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE number = :number`
	row := s.db.QueryRow(query, sql.Named("number", number))
	var p Parcel
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Получает все посылки клиента по его идентификатору
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE client = :client`
	rows, err := s.db.Query(query, sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}
	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// Обновляет статус посылки
	query := `UPDATE parcel SET status = :status WHERE number = :number`
	_, err := s.db.Exec(query, sql.Named("status", status), sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Обновляет адрес доставки посылки, если её статус 'registered'
	query := `UPDATE parcel SET address = :address WHERE number = :number AND status = :status`
	_, err := s.db.Exec(query, sql.Named("address", address), sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	return err
}

func (s ParcelStore) Delete(number int) error {
	// Удаляет посылку по её номеру, если её статус 'registered'
	query := `DELETE FROM parcel WHERE number = :number AND status = :status`
	_, err := s.db.Exec(query, sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	return err
}
