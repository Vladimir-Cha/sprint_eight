package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parsel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	req, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(req), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	res := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	p := Parcel{}
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	for rows.Next() {
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	var res []Parcel
	for _, v := range res {
		v.Number = p.Number
		v.Client = p.Client
		v.Status = p.Status
		v.Address = p.Address
		v.CreatedAt = p.CreatedAt
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE status = :status",
		sql.Named("address", address),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE status = :status",
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
