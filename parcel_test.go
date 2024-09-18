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
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "ошибка при открытии БД")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	assert.NoError(t, err, "ошибка добавления посылки")

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	fetchedParcel, err := store.Get(id)
	assert.NoError(t, err, "ошибка получения посылки")

	assert.Equal(t, parcel.Client, fetchedParcel.Client, "поле Client не совпадает")
	assert.Equal(t, parcel.Status, fetchedParcel.Status, "поле Status не совпадает")
	assert.Equal(t, parcel.Address, fetchedParcel.Address, "поле Address не совпадает")
	assert.Equal(t, parcel.CreatedAt, fetchedParcel.CreatedAt, "поле CreatedAt не совпадает")

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	assert.NoError(t, err, "ошибка удаления посылки")

	_, err = store.Get(id)
	assert.Error(t, err, "посылка с номером должна была быть удалена")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "ошибка при открытии БД")
	defer db.Close()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err, "ошибка добавления посылки")
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.NoError(t, err, "ошибка обновления адреса")

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	updatedParcel, err := store.Get(id)
	assert.NoError(t, err, "ошибка получения обновленной посылки")

	assert.Equal(t, newAddress, updatedParcel.Address, "адрес должен быть обновлен")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "ошибка при открытии БД")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	assert.NoError(t, err, "ошибка добавления посылки")

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := "Delivered"
	err = store.SetStatus(id, newStatus)
	assert.NoError(t, err, "ошибка обновления статуса")

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	updatedParcel, err := store.Get(id)
	assert.NoError(t, err, "ошибка получения обновленной посылки")

	assert.Equal(t, newStatus, updatedParcel.Status, "статус должен быть обновлен")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite3", "tracker.db")
	require.NoError(t, err, "ошибка при открытии БД")
	defer db.Close()

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
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		assert.NoError(t, err, "ошибка добавления посылки")

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	storedParcels, err := store.GetByClient(client)

	assert.NoError(t, err, "ошибка получения посылок по идентификатору клиента")
	assert.Equal(t, len(parcels), len(storedParcels), "количество полученных посылок должно совпадать с количеством добавленных")

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		originalParcel, exists := parcelMap[parcel.Number]
		assert.True(t, exists, "ожидалась посылка с номером %d в parcelMap", parcel.Number)

		// Проверяем, что значения полей полученных посылок заполнены верно
		assert.Equal(t, originalParcel.Client, parcel.Client, "поле Client не совпадает")
		assert.Equal(t, originalParcel.Status, parcel.Status, "поле Status не совпадает")
		assert.Equal(t, originalParcel.Address, parcel.Address, "поле Address не совпадает")
		assert.Equal(t, originalParcel.CreatedAt, parcel.CreatedAt, "поле CreatedAt не совпадает")
	}
}
