package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(number), nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	row, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()
	parcels := []Parcel{}
	for row.Next() {
		var p Parcel
		err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
		parcels = append(parcels, p)

	}
	return parcels, nil
}

func (s ParcelStore) GetParcelByID(parcelID int) (Parcel, error) {
	rows := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", parcelID))
	var p Parcel
	switch err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err {
	case sql.ErrNoRows:
		return p, errors.New("no rows")
	case nil:
		return p, nil
	default:
		return p, err
	}
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status =:status WHERE number =: number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status", sql.Named("number", number), sql.Named("status", "sent"))
	if err != nil {
		return err
	}
	return nil
}
