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

// getTestParcel вернет тестовую осылку
// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// проверяет добавление посылки
func TestAdd(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") //настроить подключение к БД
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	//add
	//добавить новую посылку в бд,
	par := getTestParcel()
	_, err = store.Add(par)

	//убедиться в отсутствии ошибки
	require.NoError(t, err)
	//и наличии идентификатора
	require.NotEmpty(t, par.Client)

	//проверить, что значения всех полей в полученном о бъекте совпадают со значениями полей в переменной Parcel
	//правильно ли я поняла, как организовать проверку по всем полям?
	assert.Equal(t, parcel.Client, par.Client)
	assert.Equal(t, parcel.Status, par.Status)
	assert.Equal(t, parcel.Address, par.Address)
	assert.Equal(t, parcel.CreatedAt, par.CreatedAt)

}

func TestGet(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") //настроить подключение к БД
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	//getTestParcel
	//получить только что добавленную посылку,
	_, err = store.GetByClient(parcel.Client)
	//убедиться в отсутствии ошибки
	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") //настроить подключение к БД
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	//getTestParcel
	//получить только что добавленную посылку,
	_, err = store.GetByClient(parcel.Client)
	//убедиться в отсутствии ошибки
	require.NoError(t, err)

	//delete
	//удалить добавленную посылку,
	err = store.Delete(parcel.Client)
	//убедиться в отсутствии ошибки
	require.NoError(t, err)
}

func TestSetAddress(t *testing.T) {
	//prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	//add
	//добавить новую посылку в бд, убедиться в отсутствии ошибки и наличии идентификатора
	par := getTestParcel()
	id, err := store.Add(par)
	//убедиться в отсутствии ошибки
	require.NoError(t, err)
	//и наличии идентификатора
	require.NotEmpty(t, id)

	//set address
	//обноввить адрес, убедиться в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(par.Number, newAddress)
	require.NoError(t, err)

	//check
	//получить добавленную посылку и убедиться, что статус обновился
	_, err = store.GetByClient(par.Client)
	//убедиться в отсутствии ошибки
	require.Equal(t, par.Status, err)

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	//prepare
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

	//задаем всем посылкам один и тот же идентификатора
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	//add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) //добавить новую посылку в БД
		//убедиться в отсутствии ошибки и наличии идентификатора
		if err != nil {
			require.NoError(t, err)
		}
		//обноввить идентификатор у добавленной посылки
		parcels[i].Number = int(id)
		//сохранить добавленную посылку в структуру map, чтобы ее можно было легко достать по идентификатору
		parcelMap[int(id)] = parcels[i]

		//get by client
		storedParcels, err := store.GetByClient(parcels[i].Client) //получить список посылок по идентификатору клиента, сохраненного в переменной client
		//убедиться в отсутствии ошибки
		if err != nil {
			require.NoError(t, err)
		}
		//убедиться, что количество полученных посылок совпадает с количеством добавленных
		require.Equal(t, storedParcels, err)

		//check
		//правильно ли организована проверка?
		for _, parcel := range parcelMap {
			for _, par := range parcels {
				assert.Equal(t, parcel, par)
				assert.Equal(t, parcel.Client, par.Client)
				assert.Equal(t, parcel.Status, par.Status)
				assert.Equal(t, parcel.Address, par.Address)
				assert.Equal(t, parcel.CreatedAt, par.CreatedAt)
			}
			//в parcelMap лежат добавленные посылки, ключ-идентификатор посылки,
			//значение - сама посылка
			//убедиться, что все посылки из storedParcels есть в parcelMap
			//проверить, что значения всех полей в полученном о бъекте совпадают со значениями полей в переменной Parcel

			//убедиться, что значения полей полученных посылок заполнены верно
		}
	}
}
