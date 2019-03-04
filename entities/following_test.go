package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFollowing(t *testing.T) {
	client := ConnectToRedis()
	defer FlushRedis(client)

	userLogin1 := "tavaresdong"
	userLogin2 := "king of go"

	uid1 := CreateUser(client, userLogin1, "yucdong")
	uid2 := CreateUser(client, userLogin2, "king")

	assert.NotEqual(t, "", uid1)
	assert.NotEqual(t, "", uid2)

	// The Followee posts a message
	postedMessage := "What a beautiful day today"
	_, err := CreateStatus(client, uid2, postedMessage)
	assert.Nil(t, err)

	// The Follower follows the Followee
	err = FollowUser(client, uid1, uid2)
	assert.Nil(t, err)

	messages, err := GetStatusMessage(client, uid1, "home:", 1, 30)
	assert.Nil(t, err)
	assert.Len(t, messages, 1)

	messages, err = GetStatusMessage(client, uid1, "home:", 1, 30)
	assert.Nil(t, err)
	assert.Len(t, messages, 1)

}

func TestFollowAndUnfollow(t *testing.T) {
	client := ConnectToRedis()
	defer FlushRedis(client)

	userLogin1 := "tavaresdong"
	userLogin2 := "king of go"

	uid1 := CreateUser(client, userLogin1, "yucdong")
	uid2 := CreateUser(client, userLogin2, "king")

	assert.NotEqual(t, "", uid1)
	assert.NotEqual(t, "", uid2)

	// The Followee posts a message
	postedMessage := "What a beautiful day today"
	_, err := CreateStatus(client, uid2, postedMessage)
	assert.Nil(t, err)

	// The Follower follows the Followee
	err = FollowUser(client, uid1, uid2)
	assert.Nil(t, err)

	messages, err := GetStatusMessage(client, uid1, "home:", 1, 30)
	assert.Nil(t, err)
	assert.Len(t, messages, 1)

	err = UnfollowUser(client, uid1, uid2)
	assert.Nil(t, err)

	messages, err = GetStatusMessage(client, uid1, "home:", 1, 30)
	assert.Nil(t, err)
	assert.Len(t, messages, 0)
}