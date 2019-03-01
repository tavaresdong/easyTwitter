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

func TestPostUserStatus(t *testing.T) {
	client := ConnectToRedis()
	userLogin := "tavaresdong"

	uid := CreateUser(client, userLogin, "yucdong")
	assert.NotEqual(t, "", uid)

	postedMessage := "What a beautiful day today"
	statusid, err := CreateStatus(client, uid, postedMessage)
	assert.Nil(t, err)

	data, err := FetchStatus(client, statusid)
	assert.Nil(t, err)
	assert.Equal(t, data["message"], postedMessage)
}
