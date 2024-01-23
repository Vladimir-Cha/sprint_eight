package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	result, err := s.db.Exec("INSERT INTO parcel (client, status, address) VALUES (?, ?, ?)", p.Client, p.Status, p.Address)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// Выполняем запрос к базе данных для получения данных по номеру
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = ?", number)

	// Создаем переменные для хранения данных из строки результата
	var id, client int
	var status, address string

	// Сканируем данные из строки в переменные
	err := row.Scan(&id, &client, &status, &address)
	if err != nil {
		return Parcel{}, err
	}

	// Создаем объект Parcel и заполняем его данными
	p := Parcel{
		ID:      id,
		Client:  client,
		Status:  status,
		Address: address,
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Выполняем запрос к базе данных для получения данных по клиенту
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем срез Parcel для хранения данных
	var res []Parcel

	// Итерируем по строкам результата
	for rows.Next() {
		// Создаем переменные для хранения данных из строки результата
		var id, client int
		var status, address string

		// Сканируем данные из строки в переменные
		err := rows.Scan(&id, &client, &status, &address)
		if err != nil {
			return nil, err
		}

		// Создаем объект Parcel и добавляем его к срезу
		p := Parcel{
			ID:      id,
			Client:  client,
			Status:  status,
			Address: address,
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Проверяем статус записи перед обновлением адреса
	var currentStatus string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	// Проверяем, можно ли менять адрес (если статус - "registered")
	if currentStatus != "registered" {
		return errors.New("нельзя менять адрес для посылки со статусом " + currentStatus)
	}

	// Выполняем запрос к базе данных для обновления адреса
	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// Проверяем статус записи перед удалением
	var currentStatus string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	// Проверяем, можно ли удалять запись (если статус - "registered")
	if currentStatus != "registered" {
		return errors.New("нельзя удалять посылку со статусом " + currentStatus)
	}

	// Выполняем запрос к базе данных для удаления записи
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	return err
}
