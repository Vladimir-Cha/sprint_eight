package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	addId, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, addId)

	// get
	getId, err := store.Get(addId)
	parcel.Number = getId.Number
	assert.NoError(t, err)
	assert.NotEmpty(t, parcel, getId)

	// delete
	err = store.Delete(addId)
	require.NoError(t, err)
	_, err = store.Get(addId)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	addId, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, addId)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(addId, newAddress)
	require.NoError(t, err)

	// check
	getId, err := store.Get(addId)
	require.NoError(t, err)
	require.Equal(t, newAddress, getId.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	addId, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, addId)

	// set status
	newStatus := ParcelStatusSent

	// check
	getId, err := store.Get(addId)
	require.NoError(t, err)
	require.Equal(t, newStatus, getId.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcelMap), storedParcels)
	require.NotEmpty(t, storedParcels)

	// check
	for _, parcel := range storedParcels {
		require.NotEmpty(t, parcelMap)
		require.Equal(t, len(parcelMap), len(storedParcels))
		for i := 0; i < len(parcelMap); i++ {
			require.Equal(t, parcel, parcelMap[i])
		}
	}
}
