package models

import "testing"

type testLinkedStorage struct {
	items map[itemKey]BidRecI
	users map[userKey]BidRecI
}

func newTestLinkedStorage() *testLinkedStorage {
	return &testLinkedStorage{
		items: make(map[itemKey]BidRecI),
		users: make(map[userKey]BidRecI),
	}
}

func (o *testLinkedStorage) Push(b BidRecI) {
	o.items[b.ItemId()] = b
	o.users[b.UserId()] = b
}

func TestBidTable_Push(t *testing.T) {

}
