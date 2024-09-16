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

// Добавление строки в таблицу parcel
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("insert into parcel ( client, status, address, created_at) values ( :cl,:status,:address, :creat)",
		sql.Named("cl", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("creat", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Чтение строк из таблицы parcel по заданному number, возвращается только одна строка
func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("select * from parcel where number=:num", sql.Named("num", number))

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

// Чтение строк из таблицы parcel по заданному client, может вернуться несколько строк
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var p []Parcel
	rows, err := s.db.Query("select * from parcel where client=:client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	var row Parcel
	for rows.Next() {

		err := rows.Scan(&row.Number, &row.Client, &row.Status, &row.Address, &row.CreatedAt)
		if err != nil {
			return nil, err
		}
		p = append(p, row)
	}
	return p, nil
}

// Обновление статуса в таблице parcel
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("update parcel set status=:st where number=:num",
		sql.Named("st", status),
		sql.Named("num", number))
	if err != nil {
		return err
	}
	return nil
}

// Обновление адреса в таблице parcel, изменение доступно только при статусе registered
func (s ParcelStore) SetAddress(number int, address string) error {
	var status string
	row := s.db.QueryRow("select status from parcel where number =:num",
		sql.Named("num", number))

	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if ParcelStatusRegistered != status {
		return fmt.Errorf("Wrong status!")
	}

	_, err = s.db.Exec("update parcel set address=:address where number=:num",
		sql.Named("address", address),
		sql.Named("num", number))
	if err != nil {
		return err
	}
	return nil
}

// Удаление определенной строки в таблице parcel по number, удаление доступно только при статусе registered
func (s ParcelStore) Delete(number int) error {
	var status string
	row := s.db.QueryRow("select status from parcel where number =:num",
		sql.Named("num", number))

	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("Wrong status!")
	}

	_, err = s.db.Exec("delete from parcel where number=:num", sql.Named("num", number))
	if err != nil {
		return err
	}
	return nil
}
