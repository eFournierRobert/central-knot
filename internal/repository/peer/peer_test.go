package repo_peer

import (
	"central-knot/internal/repository/db_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndGetPeer(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425

	err := Add(id, ip, port)
	assert.Nil(t, err)

	peer, err := Get(id)
	assert.Nil(t, err)
	assert.Equal(t, id, peer.ClientId)
	assert.Equal(t, ip, peer.Ip)
	assert.Equal(t, port, peer.Port)
}

func TestAddExistingPeer(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425

	err := Add(id, ip, port)
	assert.Nil(t, err)

	err = Add(id, ip, port)
	assert.NotNil(t, err)
}

func TestGetNonExistingPeer(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"

	peer, err := Get(id)
	assert.NotNil(t, err)
	assert.Empty(t, peer)
}

func TestChangeIpOnExistingPeer(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425

	err := Add(id, ip, port)
	assert.Nil(t, err)

	newIp := "10.10.12.12"
	err = ChangeIp(id, newIp)
	assert.Nil(t, err)

	peer, err := Get(id)
	assert.Nil(t, err)
	assert.Equal(t, id, peer.ClientId)
	assert.Equal(t, newIp, peer.Ip)
}

func TestChangeIpOnNoneExistingPeer(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"
	newIp := "10.10.12.12"
	err := ChangeIp(id, newIp)
	assert.NotNil(t, err)

	_, err = Get(id)
	assert.NotNil(t, err)
}
