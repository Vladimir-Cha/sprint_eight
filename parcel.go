package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const ExecTimeout = 5 * time.Second

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// отсутвует проверка на number, так как это уровень бизнес-логики
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), ExecTimeout)
	defer cancelfunc()

	res, err := s.db.ExecContext(ctx, query,
		sql.Named("client", p.Client),
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		log.Printf("Error %s when inserting row into parcel table", err)
		return -1, err
	}

	newId, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s when getting result id", err)
		return -1, err
	}
	return int(newId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// отсутствует проверка на number, так как это уровень бизнес-логики
	query := "select number, client, status, address, created_at from Parcel where number=:number"
	ctx, cancelfunc := context.WithTimeout(context.Background(), ExecTimeout)
	defer cancelfunc()

	row := s.db.QueryRowContext(ctx, query, sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Printf("No parcels with num %d", number)
		default:
			log.Printf("Error %s when getting parcel with num %d", err, number)
		}
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// отсутствует проверка на clientId, так как это уровень бизнес-логики
	// Предварительный запрос, чтобы сэкономить память при большом количестве пасылок
	ctx, cancelfunc := context.WithTimeout(context.Background(), ExecTimeout)
	defer cancelfunc()

	count, err := s.getCountByClient(ctx, client)
	if err != nil {
		log.Printf("Error %s when getting parcels by clientId %d", err, client)
		return nil, err
	}
	if count == 0 {
		return []Parcel{}, nil
	}

	query := "select number, client, status, address, created_at from Parcel where client=:client"

	rows, err := s.db.QueryContext(ctx, query, sql.Named("client", client))
	defer rows.Close()

	if err != nil {
		log.Printf("Error %s when getting parcels by clientId %d", err, client)
		return nil, err
	}

	var res []Parcel
	res = make([]Parcel, 0, count)

	//итерируемся по строкам, так как между первым и вторым запросом могло пройти время
	for rows.Next() {
		var p Parcel
		rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) getCountByClient(ctx context.Context, client int) (int64, error) {
	query := "select count(*) from Parcel where client=:client"

	row := s.db.QueryRowContext(ctx, query, sql.Named("client", client))

	var count int64
	err := row.Scan(&count)

	if err != nil {
		log.Printf("Error %s when getting count of parcels for client %d", err, client)
		return -1, err
	}

	return count, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := "update parcel set status = :status where number = :number"
	ctx, cancelfunc := context.WithTimeout(context.Background(), ExecTimeout)
	defer cancelfunc()

	res, err := s.db.ExecContext(ctx, query,
		sql.Named("number", number),
		sql.Named("status", status))

	if err != nil {
		log.Printf("Error %s when update status in parcel table", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when update status in parcel table", err)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("no rows has been updated")
		log.Printf("Error %s when update status in parcel table", err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	query := "update parcel set address = :address where number = :number and status = :status"
	ctx, cancelfunc := context.WithTimeout(context.Background(), ExecTimeout)
	defer cancelfunc()

	res, err := s.db.ExecContext(ctx, query,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		log.Printf("Error %s when update address in parcel table", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when update address in parcel table", err)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("no rows has been updated")
		log.Printf("Error %s when update address in parcel table", err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	query := "delete from parcel where number = :number and status = :status"
	ctx, cancelfunc := context.WithTimeout(context.Background(), ExecTimeout)
	defer cancelfunc()

	res, err := s.db.ExecContext(ctx, query,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		log.Printf("Error %s when delete %d from parcel table", err, number)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when delete %d from parcel table", err, number)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("no rows has been deleted")
		log.Printf("Error %s when delete %d from parcel table", err, number)
		return err
	}

	return nil
}
