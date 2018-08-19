package models

import (
	"math/rand"
	"sync"
)

type UserI interface {
	Id() UserKey
	SetName(name string)
	Name() string
	Bids() []BidRecI
}

type UserKey int32

func nextUserKey() UserKey {
	return UserKey(rand.Int31())
}

type userInfo struct {
	sync.RWMutex

	id   UserKey
	name string
}

func (o *userInfo) Id() UserKey {
	if o == nil {
		return 0
	}

	o.RLock()
	defer o.RUnlock()

	return o.id
}

func (o *userInfo) SetName(name string) {
	o.Lock()
	defer o.Unlock()

	o.name = name
}

func (o *userInfo) Name() string {
	if o == nil {
		return ""
	}

	o.RLock()
	defer o.RUnlock()

	return o.name
}

type user struct {
	userInfo

	bids sync.Map
}

func User(name string) *user {
	return &user{
		userInfo: userInfo{
			id:   nextUserKey(),
			name: name,
		},
	}
}

func (o *user) Push(bid BidRecI) {
	o.bids.LoadOrStore(bid.Id(), bid)
}

func (o *user) Bids() []BidRecI {
	o.RLock()
	defer o.RUnlock()

	result := make([]BidRecI, 0, 100)
	o.bids.Range(func(key, value interface{}) bool {
		result = append(result, value.(BidRecI))

		return true
	})

	return result
}
