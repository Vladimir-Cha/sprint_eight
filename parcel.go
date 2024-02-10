package main

import (
	"database/sql"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// добавление строки в таблицу
func (s ParcelStore) Add(p Parcel) (int, error) {
	ins, err := s.db.Exec(
		"INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, err := ins.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return int(id), nil
}

// получение посылки по номеру
func (s ParcelStore) Get(number int) (Parcel, error) {
	parcel := Parcel{}

	row := s.db.QueryRow(
		"SELECT client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number),
	)
	err := row.Scan(&parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		log.Println(err)
		return parcel, err
	}
	return parcel, nil
}

// получение всех посылок клиента
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel

	for rows.Next() {
		parcel := Parcel{}

		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		parcels = append(parcels, parcel)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return parcels, nil
}

// обновление статуса посылки
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(
		"UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// обновление адреса в таблице, возможно только если статус - registered
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(
		"UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// удаление строки из таблицы, возможно только если статус - registered
func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec(
		"DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
