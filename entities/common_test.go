package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDouble(t *testing.T) {
	val := 2
	got := Double(val)
	if got != 4 {
		t.Errorf("Value of double %d is not 4", val)
	}
}

func TestPingServer(t *testing.T) {
	ConnectToRedis()
}

func TestAcquireLock(t *testing.T) {
	client := ConnectToRedis()
	ClearLock(client, "foo")
	identifier := AcquireLock(client, "foo", 5)
	assert.NotEqual(t, "", identifier)

	newId := AcquireLock(client, "foo", 15)
	assert.Equal(t, "", newId)
}

func TestReleaseLock(t *testing.T) {
	client := ConnectToRedis()
	ClearLock(client, "foo")

	identifier := AcquireLock(client, "foo", 15)
	assert.NotEqual(t, "", identifier)
	err := ReleaseLock(client, "foo", identifier)

	assert.Nil(t, err)
}

func TestReleaseLockFail(t *testing.T) {
	client := ConnectToRedis()
	ClearLock(client, "foo")

	AcquireLock(client, "foo", 15)
	err := ReleaseLock(client, "foo", "bar")

	assert.NotNil(t, err)
}


func TestAcquireLockAndWait(t *testing.T) {
	client := ConnectToRedis()
	ClearLock(client, "foo")

	identifier := AcquireLock(client, "foo", 5)
	assert.NotEqual(t, "", identifier)

	newId := AcquireLockWithTimeout(client, "foo", 15, 5)
	assert.NotEqual(t, "", newId)
}

func TestAcquireLockAndWaitTimeout(t *testing.T) {
	client := ConnectToRedis()
	cleared := ClearLock(client, "foo")
	assert.True(t, cleared)

	identifier := AcquireLock(client, "foo", 15)
	assert.NotEqual(t, "", identifier)

	newId := AcquireLockWithTimeout(client, "foo", 5, 5)
	assert.Equal(t, "", newId)
}