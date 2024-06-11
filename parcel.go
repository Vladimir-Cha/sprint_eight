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
	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (?, ?, ?, ?, ?)",
		p.Number, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Вернуть идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = ?", number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем пустой срез Parcel, в который будем собирать результаты запроса
	var res []Parcel

	// Итерируемся по результатам запроса
	for rows.Next() {
		var p Parcel
		// Сканируем значения из строк запроса в структуру Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		// Добавляем полученную посылку в срез
		res = append(res, p)
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращаем результаты запроса и ошибку (если есть)
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// Выполняем запрос на обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return err
	}

	// Возвращаем ошибку (если есть)
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Проверяем статус посылки
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number and status = :status",
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("number", number),
		sql.Named("address", address))

	return err
}

func (s ParcelStore) Delete(number int) error {
	// Проверяем статус посылки
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number and status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	return err
}
