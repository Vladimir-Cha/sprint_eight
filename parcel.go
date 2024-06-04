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

func (s ParcelStore) Add(p Parcel) (int64, error) {
	res, err := s.db.Exec("INSERT INTO parcels  (parcelName, parcelStatus) VALUES (:parcelName, :prcelStatus)",
		sql.Named("name", p.Address),
		sql.Named("status", p.Status))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	fmt.Println(res.LastInsertId())
	fmt.Println(res.RowsAffected())

	return res.LastInsertId()
}

// 2. Чтение строки из БД по заданному client - получение списка посылок клиента
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELCT * FROM parcels WHERE client = :client", sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	//заполнить объект Parcel данными из таблицы
	parcels := []Parcel{}
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
		parcels = append(parcels, p)

	}
	return parcels, nil
}

func (s ParcelStore) GetParcelByID(parcelID int) (Parcel, error) {
	rows, err := s.db.Query("SELCT * FROM parcels WHERE id = :id", sql.Named("id", parcelID))
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var p Parcel
	for rows.Next() {

		err := rows.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}

	}
	return p, nil
}

// 4. Изменение статуса посылки
func (s ParcelStore) SetStatus(number int, status string) error {
	//реализовать обновление статуса в таблице Parcel
	_, err := s.db.Exec("UPDATE parcel SET status =:status WHERE number =: number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// 4. Изменение адреса доставки
func (s ParcelStore) SetAddress(number int, address string) error {
	//реализовать обновление адреса в таблице parcel
	//менять можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address =:address WHERE number =: number AND status =: status", sql.Named("address", address), sql.Named("number", number), sql.Named("status", "registered"))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// 5. Удаление посылки
func (s ParcelStore) Delete(number int) error {
	//реализовать удаление строки из таблицы parcel
	//удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number =: number AND status =: status", sql.Named("number", number), sql.Named("status", "registered"))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
