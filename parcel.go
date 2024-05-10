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
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}

	lastID, _ := res.LastInsertId()
	// верните идентификатор последней добавленной записи
	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	p := Parcel{}

	row := s.db.QueryRow("SELECT address, client, created_at, number, status FROM parcel WHERE number = :number",
		sql.Named("number", number),
	)
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы

	err := row.Scan(&p.Address, &p.Client, &p.CreatedAt, &p.Number, &p.Status)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query("SELECT address, client, created_at, number, status FROM parcel WHERE client = :client",
		sql.Named("client", client),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// здесь из таблицы может вернуться несколько строк

	var res []Parcel

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Address, &p.Client, &p.CreatedAt, &p.Number, &p.Status)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	// заполните срез Parcel данными из таблицы
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return err
	}
	// менять адрес можно только если значение статуса registered

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return err
	}
	// удалять строку можно только если значение статуса registered

	return nil
}
