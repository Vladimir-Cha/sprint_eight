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
	// добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel ( client, status, address, created_at) VALUES ( :client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		werr := fmt.Errorf("ошибка добавления посылки. ошибка: %w", err)
		return 0, werr
	}
	// возвращаем идентификатор последней добавленной записи
	insId, err := res.LastInsertId()
	if err != nil {
		werr := fmt.Errorf("ошибка чтения идентификатора добавленной посылки. ошибка: %w", err)
		return int(insId), werr
	}

	return int(insId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// чтение строки по заданному number
	// из таблицы должна вернуться только одна строка
	pRow := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	// заполнение объекта Parcel данными из таблицы
	p := Parcel{}
	err := pRow.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		werr := fmt.Errorf("ошибка при чтении посылки по номеру %d. ошибка: %w", number, err)
		return p, werr
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// чтение строк из таблицы parcel по заданному client
	// из таблицы может вернуться несколько строк
	pRows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))

	// заполняем срез Parcel данными из таблицы
	var res []Parcel

	if err != nil {
		werr := fmt.Errorf("ошибка при чтении посылок клиента %d. ошибка: %w", client, err)
		return res, werr
	}
	defer pRows.Close()

	for pRows.Next() {
		var parcel Parcel

		err := pRows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			werr := fmt.Errorf("ошибка при чтении посылок клиента %d. ошибка: %w", client, err)
			return res, werr
		}
		res = append(res, parcel)
	}

	if err := pRows.Err(); err != nil {
		werr := fmt.Errorf("ошибка при чтении посылок клиента %d. ошибка: %w", client, err)
		return res, werr
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		return fmt.Errorf("невозможно обновить статус. ошибка: %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("невозможно обновить адрес. ошибка: %w", err)
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("невозможно обновить адрес. недопустимый статус посылки. требуемый статус: %s, текущий статус: %s", ParcelStatusRegistered, parcel.Status)

	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number),
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении адреса. ошибка: %w ", err)

	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("невозможно удалить посылку. ошибка: %w ", err)
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("невозможно удалить посылку. недопустимый статус посылки. требуемый статус: %s, текущий статус: %s", ParcelStatusRegistered, parcel.Status)

	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number",
		sql.Named("number", number),
	)

	if err != nil {
		return fmt.Errorf("ошибка при удалении посылки. ошибка: %w ", err)

	}
	return nil
}
