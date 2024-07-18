// тут мы работаем с конкретной посылкой
package main

import (
	"database/sql"
)

// Итак, поле db типа *sql.DB - это указатель на соединение с базой данных.
type ParcelStore struct {
	db *sql.DB // это переменная, которая хранит адрес объекта sql.DB в памяти компьютера. Объект sql.DB представляет собой соединение с базой данных SQLite.
	// Когда вы выполняете операцию с базой данных с помощью объекта db, он отправляет соответствующий запрос на сервер базы данных. Сервер базы данных выполняет запрос и возвращает результаты обратно объекту db. Затем объект db преобразует результаты в формат, который может быть использован приложением.
	// Важно отметить, что объект db не хранит данные из базы данных в памяти компьютера. Он просто предоставляет канал связи для отправки запросов и получения результатов.
}

// NewParcelStore() получает подключение к базе данных и возвращает объект структуры ParcelStore, который содержит в себе это подключение, а также набор методов для работы с посылками в базе данных.
// Таким образом, объект структуры ParcelStore содержит в себе ссылку на то же самое подключение к базе данных, которое было передано функции NewParcelStore(). Однако объект структуры ParcelStore также предоставляет набор методов для работы с посылками в базе данных, таких как Add(), Get() и тд.
// Первое db - это параметр функции NewParcelStore(), который представляет собой подключение к базе данных типа *sql.DB.
// Второе db - это поле структуры ParcelStore (которое также имеет тип *sql.DB).
// Структура ParcelStore имеет поле db типа *sql.DB. Когда вы создаете новую структуру ParcelStore, вы должны инициализировать это поле подключением к базе данных. Вы делаете это, передав подключение к базе данных в качестве аргумента функции NewParcelStore().
// Внутри функции NewParcelStore() вы создаете новый объект структуры ParcelStore и инициализируете его поле db переданным подключением к базе данных.
// Таким образом, запись {db: db} означает, что поле db структуры ParcelStore инициализируется значением параметра db, который является адресом в памяти, где хранятся параметры подключения к базе данных.
// Функция NewParcelStore() не должна использоваться в каждой функции, потому что она используется только для создания нового объекта хранилища посылок. После того, как объект хранилища создан, его можно использовать для работы с посылками в базе данных, передавая его в другие функции в качестве аргумента.
// Мы передаем хранилище посылок в качестве аргумента функциям Add() и Get(). Это позволяет этим функциям использовать подключение к базе данных, хранящееся в хранилище посылок, для выполнения своих операций.
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	stmt, err := s.db.Prepare("INSERT INTO parcel (number, client, status, address, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(p.Number, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	// return 0, nil - не понял, что это
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	stmt, err := s.db.Prepare("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?")
	if err != nil {
		return Parcel{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(number)
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	stmt, err := s.db.Prepare("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(client)
	if err != nil {
		return nil, err
	}

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	stmt, err := s.db.Prepare("UPDATE parcel SET status = ? WHERE number = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(status, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	stmt, err := s.db.Prepare("UPDATE parcel SET address = ? WHERE number = ? AND status = 'registered'")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(address, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	stmt, err := s.db.Prepare("DELETE FROM parcel WHERE number = ? AND status = 'registered'")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(number)
	if err != nil {
		return err
	}

	return nil
}
