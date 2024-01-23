package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := tracker.db()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.AddParcel(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get
	retrievedParcel, err := store.GetParcel(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, retrievedParcel.Client)
	require.Equal(t, parcel.Status, retrievedParcel.Status)
	require.Equal(t, parcel.Address, retrievedParcel.Address)
	// Проверьте, что значения времени совпадают с некоторой разумной точностью
	require.WithinDuration(t, time.Now().UTC(), retrievedParcel.CreatedAt, time.Second)

	// delete
	err = store.DeleteParcel(id)
	require.NoError(t, err)

	// Проверьте, что посылку больше нельзя получить из БД
	_, err = store.GetParcel(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := tracker.db() // настройте подключение к БД
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.AddParcel(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	retrievedParcel, err := store.GetParcel(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, retrievedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := tracker.db() // настройте подключение к БД
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.AddParcel(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// set status
	newStatus := ParcelStatusDelivered // установите новый статус
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	retrievedParcel, err := store.GetParcel(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, retrievedParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := tracker.db() // настройте подключение к БД
	require.NoError(t, err)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.AddParcel(parcels[i])
		require.NoError(t, err)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetParcelsByClient(client)
	require.NoError(t, err)

	// check
	require.Equal(t, len(parcels), len(storedParcels))

	for _, storedParcel := range storedParcels {
		// проверяем, что посылка существует в parcelMap
		originalParcel, ok := parcelMap[storedParcel.Number]
		require.True(t, ok)

		// убедимся, что значения полей полученных посылок заполнены верно
		require.Equal(t, originalParcel.Client, storedParcel.Client)
		require.Equal(t, originalParcel.Status, storedParcel.Status)
		require.Equal(t, originalParcel.Address, storedParcel.Address)
	}
}
