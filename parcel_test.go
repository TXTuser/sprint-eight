package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// источник псевдослучайных чисел
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// тестовая посылка
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
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	// ADD
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// GET
	storedParcel, err := store.Get(id)
	require.NoError(t, err)

	parcel.Number = id
	require.Equal(t, parcel, storedParcel)

	// DELETE
	err = store.Delete(id)
	require.NoError(t, err)

	// проверяем что записи больше нет
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	storedParcel, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, newAddress, storedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	storedParcel, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, ParcelStatusSent, storedParcel.Status)
}

// TestGetByClient проверяет получение посылок по клиенту
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)

	for i := 0; i < len(parcels); i++ {

		parcels[i].Client = client

		id, err := store.Add(parcels[i])
		require.NoError(t, err)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	require.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {

		expectedParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok)

		require.Equal(t, expectedParcel, parcel)
	}
}
