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
	query := `INSERT INTO parcels (client, status, address, created_at) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :id",
		sql.Named("id", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT number, client, status, address, created_at FROM parcels WHERE client = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := `UPDATE parcels SET status = ? WHERE number = ?`
	_, err := s.db.Exec(query, status, number)
	return err
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
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
