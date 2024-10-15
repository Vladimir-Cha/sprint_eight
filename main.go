package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

type Parcel struct {
	Number    int
	Client    int
	Status    stringAddress   string
	CreatedAt string
}

type ParcelStore interface {
	Add(parcel Parcel) (int, error)
	Get(number int) (Parcel, error)
	GetByClient(client int) ([]Parcel, error)
	SetStatus(number int, status string) error
	SetAddress(number int, address string) error
	Delete(number int) error
}

type SQLiteParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return &SQLiteParcelStore{db: db}
}

func (s *SQLiteParcelStore) Add(parcel Parcel) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO parcels (client, status, address, created_at) VALUES (?, ?, ?, ?)",
		parcel.Client, parcel.Status, parcel.Address, parcel.CreatedAt,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SQLiteParcelStore) Get(number int) (Parcel, error) {
	var parcel Parcel
	err := s.db.QueryRow(
		"SELECT number, client, status, address, created_at FROM parcels WHERE number = ?",
		number,
	).Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return parcel, nil
}

func (s *SQLiteParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcels WHERE client = ?",
		client,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var parcel Parcel
		if err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt); err != nil {
			return nil, err	}
		parcels = append(parcels, parcel)
	}
	return parcels, nil
}

func (s *SQLiteParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcels SET status = ? WHERE number = ?", status, number)
	return err
}

func (s *SQLiteParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcels SET address = ? WHERE number = ?", address, number)
	return err
}

func (s *SQLiteParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcels WHERE number = ?", number)
	return err
}

type ParcelService struct {
	store ParcelStore
}

func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

func (s ParcelService) Register(client int, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err}

	parcel.Number = id

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)

	return parcel, nil
}

func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println()

	return nil
}

func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	case ParcelStatusDelivered:
		return nil}

	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)

	return s.store.SetStatus(number, nextStatus)
}

func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

func main() {
	// Открытие соединения с базой данных
	db, err := sql.Open("sqlite", "parcels.db")
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	// Создание таблицы, если она не существует_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS parcels (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT	)
	`)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
		return
	}

	// Создание объекта ParcelStore
	store := NewParcelStore(db)
	service := NewParcelService(store)

	// Пример использования сервиса
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Изменение статуса
	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Попытка удаления отправленной посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Регистрация новой посылки
	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Удаление новой посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
