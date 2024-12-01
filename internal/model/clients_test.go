package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestICanAddAndRemoveClients(t *testing.T) {
	clients := NewClients()
	cl1_1 := Client{
		Event:        "event1",
		Uuid:         "uuid1",
		SubscribedAt: time.Time{},
	}
	cl2_2 := Client{
		Event:        "event2",
		Uuid:         "uuid2",
		SubscribedAt: time.Time{},
	}
	assert.NoError(t, clients.Add(&cl1_1))
	assert.NoError(t, clients.Add(&cl2_2))
	assert.Equal(t, &cl1_1, clients["event1"][0])
	assert.Equal(t, &cl2_2, clients["event2"][0])
	assert.NoError(t, clients.Remove(&cl1_1))
	assert.Equal(t, 0, len(clients["event1"]))
	cl1_3 := Client{
		Event:        "event1",
		Uuid:         "uuid3",
		SubscribedAt: time.Time{},
	}
	// removing all of the same event, then add again
	assert.NoError(t, clients.Add(&cl1_3))
	assert.Equal(t, &cl1_3, clients["event1"][0])
	cl2_4 := Client{
		Event:        "event2",
		Uuid:         "uuid4",
		SubscribedAt: time.Time{},
	}
	assert.NoError(t, clients.Add(&cl2_4))
	assert.Equal(t, 2, len(clients["event2"]))
	assert.NoError(t, clients.Remove(&cl1_3))
	// trying out of bound test
	assert.NoError(t, clients.Remove(&cl1_3))
	assert.Equal(t, 0, len(clients["event1"]))
	assert.NoError(t, clients.Remove(&cl2_2))
	assert.NoError(t, clients.Remove(&cl2_4))
	assert.Equal(t, 0, len(clients["event2"]))
	// trying out of bound tests
	assert.NoError(t, clients.Remove(&cl2_2))
	assert.NoError(t, clients.Remove(&cl2_4))
	assert.Equal(t, 0, len(clients["event2"]))
}

func TestClientsNilParamsAndRemoveOnNonExistingEvent(t *testing.T) {
	clients := NewClients()
	assert.ErrorContains(t, clients.Add(nil), ErrCouldNotAddClient.Error())
	assert.ErrorContains(t, clients.Add(nil), ErrNilParameter.Error())
	assert.ErrorContains(t, clients.Remove(nil), ErrCouldNotRemoveClient.Error())
	assert.ErrorContains(t, clients.Remove(nil), ErrNilParameter.Error())
	cl1_1 := Client{
		Event:        "event_doesnt_exist",
		Uuid:         "uuid1",
		SubscribedAt: time.Time{},
	}
	assert.NoError(t, clients.Remove(&cl1_1))
	assert.NoError(t, clients.Remove(&Client{}))
}
