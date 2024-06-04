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
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
func TestAdd(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	par := getTestParcel()
	_, err = store.Add(par)
	require.NoError(t, err)
	require.NotEmpty(t, par.Client)
	assert.Equal(t, parcel.Client, par.Client)
	assert.Equal(t, parcel.Status, par.Status)
	assert.Equal(t, parcel.Address, par.Address)
	assert.Equal(t, parcel.CreatedAt, par.CreatedAt)

}

func TestGet(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	_, err = store.GetByClient(parcel.Client)
	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	_, err = store.GetByClient(parcel.Client)
	require.NoError(t, err)
	err = store.Delete(parcel.Client)
	require.NoError(t, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	par := getTestParcel()
	id, err := store.Add(par)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	newAddress := "new test address"
	err = store.SetAddress(par.Number, newAddress)
	require.NoError(t, err)
	_, err = store.GetByClient(par.Client)
	require.Equal(t, par.Status, err)

}
func TestGetByClient(t *testing.T) {
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

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		if err != nil {
			require.NoError(t, err)
		}
		parcels[i].Number = int(id)
		parcelMap[int(id)] = parcels[i]

		storedParcels, err := store.GetByClient(parcels[i].Client)
		if err != nil {
			require.NoError(t, err)
		}
		require.Equal(t, storedParcels, err)
		for _, parcel := range parcelMap {
			for _, par := range parcels {
				assert.Equal(t, parcel, par)
				assert.Equal(t, parcel.Client, par.Client)
				assert.Equal(t, parcel.Status, par.Status)
				assert.Equal(t, parcel.Address, par.Address)
				assert.Equal(t, parcel.CreatedAt, par.CreatedAt)
			}
		}
	}
}
