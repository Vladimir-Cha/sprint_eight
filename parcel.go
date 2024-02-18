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
	// реализуйте добавление строки в таблицу parcel,
	//используйте данные из переменной p
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return p, nil
	}
	defer db.Close()
	row := db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	_ = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return nil, nil
	}
	defer db.Close()
	row, err := db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, nil
	}
	defer row.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for row.Next() {
		pars := Parcel{}

		err := row.Scan(&pars.Number, &pars.Client, &pars.Status, &pars.Address, &pars.CreatedAt)
		if err != nil {
			return nil, nil
		}
		res = append(res, pars)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	return nil
}
