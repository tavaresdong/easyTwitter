package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCrud(t *testing.T) {
	client := ConnectToRedis()
	userLogin := "tavaresdong"
	DeleteUser(client, userLogin)

	uid := CreateUser(client, userLogin, "yucdong")
	assert.NotEqual(t, "", uid)

	info := FetchUser(client, userLogin)

	assert.NotNil(t, info)
	fmt.Println(info["name"])
	fmt.Println(info["signup"])

	err := DeleteUser(client, userLogin)
	assert.Nil(t, err)
}
