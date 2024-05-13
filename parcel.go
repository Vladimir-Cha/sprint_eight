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

	res, err := s.db.Exec("INSERT INTO parcel (Client, Status, Address, Created_at) VALUES (:Client, :Status, :Address, :Created_at)",
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("Created_at", p.CreatedAt))
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

	p := Parcel{}

	row := s.db.QueryRow("SELECT * FROM parcel WHERE Number = :Number", sql.Named("Number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	rows, err := s.db.Query("SELECT * FROM parcel WHERE Client = :Client", sql.Named("Client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		parcel := Parcel{}
		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, parcel)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET Status = :Status WHERE Number = :Number",
		sql.Named("Number", number),
		sql.Named("Status", status))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	_, err := s.db.Exec("UPDATE parcel SET Address = :Address WHERE Number = :Number and Status = :Status",
		sql.Named("Number", number),
		sql.Named("Address", address),
		sql.Named("Status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE Number = :Number and Status = :Status",
		sql.Named("Number", number),
		sql.Named("Status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}
