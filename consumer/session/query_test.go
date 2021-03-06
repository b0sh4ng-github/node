/*
 * Copyright (C) 2020 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package session

import (
	"math/big"
	"testing"
	"time"

	"github.com/mysteriumnetwork/node/identity"
	session_node "github.com/mysteriumnetwork/node/session"
	"github.com/stretchr/testify/assert"
)

func TestSessionQuery_FetchSessions(t *testing.T) {
	// given
	session1Expected := History{
		SessionID: session_node.ID("session1"),
		Started:   time.Date(2020, 6, 17, 0, 0, 1, 0, time.UTC),
	}
	session2Expected := History{
		SessionID: session_node.ID("session2"),
		Started:   time.Date(2020, 6, 17, 0, 0, 2, 0, time.UTC),
	}
	storage, storageCleanup := newStorageWithSessions(session1Expected, session2Expected)
	defer storageCleanup()

	// when
	query := NewQuery().FetchSessions()
	err := storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(t, []History{session2Expected, session1Expected}, query.Sessions)
}

func TestSessionQuery_FilterDirection(t *testing.T) {
	// given
	sessionExpected := History{
		SessionID: session_node.ID("session1"),
		Direction: "Provided",
	}
	storage, storageCleanup := newStorageWithSessions(sessionExpected)
	defer storageCleanup()

	// when
	query := NewQuery().FetchSessions()
	err := storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(t, []History{sessionExpected}, query.Sessions)

	// when
	query = NewQuery().FetchSessions().FilterDirection(DirectionConsumed)
	err = storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(t, []History{}, query.Sessions)
}

func TestSessionQuery_FetchStats(t *testing.T) {
	// given
	sessionExpected := History{
		SessionID:    session_node.ID("session1"),
		Direction:    "Provided",
		ConsumerID:   identity.FromAddress("consumer1"),
		DataSent:     1234,
		DataReceived: 123,
		Tokens:       big.NewInt(12),
		Started:      time.Date(2020, 6, 17, 10, 11, 12, 0, time.UTC),
		Updated:      time.Date(2020, 6, 17, 10, 11, 32, 0, time.UTC),
		Status:       "New",
	}
	storage, storageCleanup := newStorageWithSessions(sessionExpected)
	defer storageCleanup()

	// when
	query := NewQuery().FetchStats()
	err := storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(
		t,
		Stats{
			Count: 1,
			ConsumerCounts: map[identity.Identity]int{
				identity.FromAddress("consumer1"): 1,
			},
			SumDataSent:     1234,
			SumDataReceived: 123,
			SumTokens:       big.NewInt(12),
			SumDuration:     20 * time.Second,
		},
		query.Stats,
	)

	// when
	query = NewQuery().FilterDirection(DirectionConsumed).FetchStats()
	err = storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(t, NewStats(), query.Stats)
}

func TestSessionQuery_FetchStatsByDay(t *testing.T) {
	// given
	sessionExpected := History{
		SessionID:    session_node.ID("session1"),
		Direction:    "Provided",
		ConsumerID:   identity.FromAddress("consumer1"),
		DataSent:     1234,
		DataReceived: 123,
		Tokens:       big.NewInt(12),
		Started:      time.Date(2020, 6, 17, 10, 11, 12, 0, time.UTC),
		Updated:      time.Date(2020, 6, 17, 10, 11, 32, 0, time.UTC),
		Status:       "New",
	}
	storage, storageCleanup := newStorageWithSessions(sessionExpected)
	defer storageCleanup()

	// when
	query := NewQuery().
		FilterFrom(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)).
		FilterTo(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)).
		FetchStatsByDay()
	err := storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(
		t,
		map[time.Time]Stats{
			time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC): NewStats(),
		},
		query.StatsByDay,
	)

	// when
	query = NewQuery().
		FilterFrom(time.Date(2020, 6, 17, 0, 0, 0, 0, time.UTC)).
		FilterTo(time.Date(2020, 6, 18, 0, 0, 0, 0, time.UTC)).
		FetchStatsByDay()
	err = storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(
		t,
		map[time.Time]Stats{
			time.Date(2020, 6, 17, 0, 0, 0, 0, time.UTC): {
				Count: 1,
				ConsumerCounts: map[identity.Identity]int{
					identity.FromAddress("consumer1"): 1,
				},
				SumDataSent:     1234,
				SumDataReceived: 123,
				SumTokens:       big.NewInt(12),
				SumDuration:     20 * time.Second,
			},
			time.Date(2020, 6, 18, 0, 0, 0, 0, time.UTC): NewStats(),
		},
		query.StatsByDay,
	)

	// when
	query = NewQuery().
		FilterFrom(time.Date(2020, 6, 17, 0, 0, 0, 0, time.UTC)).
		FilterTo(time.Date(2020, 6, 18, 0, 0, 0, 0, time.UTC)).
		FilterDirection(DirectionConsumed).
		FetchStatsByDay()
	err = storage.Query(query)
	// then
	assert.Nil(t, err)
	assert.Equal(
		t,
		map[time.Time]Stats{
			time.Date(2020, 6, 17, 0, 0, 0, 0, time.UTC): NewStats(),
			time.Date(2020, 6, 18, 0, 0, 0, 0, time.UTC): NewStats(),
		},
		query.StatsByDay,
	)
	return
}
