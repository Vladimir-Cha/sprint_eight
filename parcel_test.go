package main

import (
	"database/sql"
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

	// Изменяем статус посылки
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// Проверяем новый статус
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, storedParcel.Status)

	// Пытаемся удалить посылку
	err = store.Delete(id)
	require.Error(t, err) // ожидаем ошибку, т.к. статус не "registered"

	// Возвращаем статус в "registered", чтобы можно было удалить
	err = store.SetStatus(id, ParcelStatusRegistered)
	require.NoError(t, err)

	// Теперь удаляем посылку
	err = store.Delete(id)
	require.NoError(t, err)
}
