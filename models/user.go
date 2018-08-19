package models

import (
	"sync"
	"math/rand"
)

type UserI interface {
	Id() userKey
	SetName(name string)
	Name() string
}

type userKey int32

func nextUserKey() userKey {
	return userKey(rand.Int31())
}

type userInfo struct {
	sync.RWMutex

	id   userKey
	name string
}

func (o *userInfo) Id() userKey {
	o.RLock()
	defer o.Unlock()

	return o.id
}

func (o *userInfo) SetName(name string) {
	o.Lock()
	defer o.Unlock()

	o.name = name
}

func (o *userInfo) Name() string {
	o.RLock()
	defer o.Unlock()

	return o.name
}

type user struct {
	userInfo

	bids sync.Map
}

func User(name string) *user {
	return &user {
		userInfo: userInfo{
			id: nextUserKey(),
			name: name,
		},
	}
}

func (o *user) Push(bid BidRecI) {
	o.bids.LoadOrStore(bid.Id(), bid)
}