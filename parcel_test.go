package main

import (
	"database/sql"
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
	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)                     // проверка ошибки

	defer db.Close()

	store := NewParcelStore(db) // создание тестовой посылки
	parcel := getTestParcel()   // получение тестовой посылки

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	checkParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, parcel, checkParcel)
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(parcel.Number)
	require.NoError(t, err)

	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(parcel.Number)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)                 // создание тестовой посылки

	defer db.Close()

	parcel := getTestParcel() // получение тестовой посылки
	require.NoError(t, err)   // проверка ошибки

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	checkParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, newAddress, checkParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)                 // создание тестовой посылки

	defer db.Close()

	parcel := getTestParcel() // получение тестовой посылки
	require.NoError(t, err)   // проверка ошибки

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)
	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(parcel.Number, ParcelStatusDelivered)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	checkParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusDelivered, checkParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	require.NoError(t, err)                     // проверка ошибки

	defer db.Close()
	store := NewParcelStore(db) // создание тестовой посылки

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{} // создание map для хранения посылок
	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ { // добавляем посылки в БД
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		require.NoError(t, err)
		require.NotEmpty(t, id)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	require.NoError(t, err)                         // убедитесь в отсутствии ошибки
	require.Len(t, storedParcels, len(parcels))     // убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		require.Contains(t, parcelMap, parcel.Number)      // каждая посылка из полученного списка есть в map
		require.Equal(t, parcelMap[parcel.Number], parcel) // убедитесь, что значения полей полученных посылок заполнены верно
	}
}
