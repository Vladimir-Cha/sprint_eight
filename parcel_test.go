package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test address",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// Подготовка
	db, err := sql.Open("sqlite3", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавление
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// Получение
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, storedParcel.Client)
	require.Equal(t, parcel.Status, storedParcel.Status)
	require.Equal(t, parcel.Address, storedParcel.Address)

	// Удаление
	err = store.Delete(id)
	require.NoError(t, err)

	// Проверка, что посылка удалена
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// Подготовка
	db, err := sql.Open("sqlite3", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавление
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// Обновление адреса
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// Проверка
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, storedParcel.Address)

	// Очистка
	err = store.Delete(id)
	require.NoError(t, err)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// Подготовка
	db, err := sql.Open("sqlite3", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавляем посылку
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// достмвляем
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// проверяем статус
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, storedParcel.Status)

	// удвляем
	err = store.Delete(id)
	require.Error(t, err) // ожидаем ошибку, т.к. статус не registered

	// ставим registred чтобы удалить
	err = store.SetStatus(id, ParcelStatusRegistered)
	require.NoError(t, err)

	// точно удаляем
	err = store.Delete(id)
	require.NoError(t, err)
}

func TestGetByClient(t *testing.T) {
	// Подключение к бд
	db, err := sql.Open("sqlite3", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	// Очистка таблицы перед тестом
	_, err = db.Exec("DELETE FROM parcel")
	require.NoError(t, err)

	// тестовые
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	clientID := rand.Intn(100500) // рандомный ID клиента
	for i := range parcels {
		parcels[i].Client = clientID // присваиваем этот id всем поссылкам
	}

	// Добавление посылок в бд и заполнение
	parcelMap := map[int]Parcel{}
	for _, parcel := range parcels {
		id, err := store.Add(parcel)
		require.NoError(t, err)
		parcelMap[id] = parcel
	}

	// GetByClient
	returnedParcels, err := store.GetByClient(clientID)
	require.NoError(t, err)

	// Проверка кол-ва записей
	require.Len(t, returnedParcels, len(parcels))

	// Проверка на соответсвие
	for _, returnedParcel := range returnedParcels {
		expectedParcel, exists := parcelMap[returnedParcel.Number]
		require.True(t, exists, "получена неизвестная посылка")
		require.Equal(t, expectedParcel.Client, returnedParcel.Client)
		require.Equal(t, expectedParcel.Address, returnedParcel.Address)
		require.Equal(t, expectedParcel.Status, returnedParcel.Status)
		require.Equal(t, expectedParcel.CreatedAt, returnedParcel.CreatedAt)
	}
}
