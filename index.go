package main

import (
	"sync"

	"github.com/inszva/gods/maps/treemap"
)

type UserRecord struct {
	*DB_User
	ref *User
}

func (userRecord *UserRecord) GetRef() *User {
	if userRecord.ref != nil { //likely
		return userRecord.ref
	}
	userRecord.ref = &User{*userRecord.DB_User}
	return userRecord.ref
}

type PaymentRecord struct {
	*DB_Payment
	ref *Payment
}

func (paymentRecord *PaymentRecord) GetRef() *Payment {
	if paymentRecord.ref != nil { //likely
		return paymentRecord.ref
	}
	paymentRecord.ref = &Payment{*paymentRecord.DB_Payment}
	return paymentRecord.ref
}

type ItemRecord struct {
	*DB_Item
	ref *Item
}

func (itemRecord *ItemRecord) GetRef() *Item {
	if itemRecord.ref != nil { //likely
		return itemRecord.ref
	}
	itemRecord.ref = &Item{*itemRecord.DB_Item}
	return itemRecord.ref
}

type UAVRecord struct {
	*DB_UAV
	ref *UAV
}

func (uavRecord *UAVRecord) GetRef() *UAV {
	if uavRecord.ref != nil { //likely
		return uavRecord.ref
	}
	uavRecord.ref = &UAV{*uavRecord.DB_UAV, make(chan int, 1)}
	return uavRecord.ref
}

type UserIndex struct {
	lock sync.RWMutex
	tree *treemap.Map
}

type UserIdTimeUnion struct {
	userId int
	time   int
}

type PaymentUserIdTimeIndex struct {
	lock sync.RWMutex
	tree *treemap.Map
}

type PaymentIndex struct {
	lock sync.RWMutex
	tree *treemap.Map
}

type ItemIndex struct {
	lock sync.RWMutex
	tree *treemap.Map
}

type UAVIndex struct {
	lock sync.RWMutex
	tree *treemap.Map
}

var (
	user_id_index = UserIndex{
		tree: treemap.NewWithIntComparator(),
	}

	payment_user_id_time_index = PaymentUserIdTimeIndex{
		tree: treemap.NewWith(func(a, b interface{}) int {
			union1, union2 := a.(UserIdTimeUnion), b.(UserIdTimeUnion)
			if union1.userId > union2.userId {
				return 1
			}
			if union1.userId < union2.userId {
				return -1
			}
			if union1.time > union2.time {
				return 1
			}
			if union1.time < union2.time {
				return -1
			}
			return 0
		}),
	}

	payment_index = PaymentIndex{
		tree: treemap.NewWithIntComparator(),
	}

	item_index = ItemIndex{
		tree: treemap.NewWithIntComparator(),
	}

	uav_index = UAVIndex{
		tree: treemap.NewWithIntComparator(),
	}
)

func (user_id_index *UserIndex) getUserById(id int) *User {
	user_id_index.lock.RLock()
	defer user_id_index.lock.RUnlock()

	_userRecord, ok := user_id_index.tree.Get(id)
	if !ok {
		return nil
	}
	userRecord, ok := _userRecord.(UserRecord)
	if !ok {
		return nil
	}
	return userRecord.GetRef()
}

func (payment_user_id_time_index *PaymentUserIdTimeIndex) insertPayment(payment Payment) {
	payment_user_id_time_index.lock.Lock()
	defer payment_user_id_time_index.lock.Unlock()

	paymentTable = append(paymentTable, payment.DB_Payment)
	payment_user_id_time_index.tree.Put(UserIdTimeUnion{
		payment.Payment_user_id, payment.Payment_time,
	}, PaymentRecord{
		DB_Payment: &paymentTable[len(paymentTable)-1],
	})
}

func (payment_index *PaymentIndex) getPayment(id int) *Payment {
	key, ok := payment_index.tree.Get(id)
	if !ok {
		return nil
	}
	_paymentRecord, ok := payment_user_id_time_index.tree.Get(key.(UserIdTimeUnion))
	if !ok {
		panic("Help index is not consistence with primary index.")
	}
	paymentRecord := _paymentRecord.(PaymentRecord)
	return paymentRecord.GetRef()
}

func (payment_user_id_time_index *PaymentUserIdTimeIndex) getUserLastPayments(userId, start, limit int) (payments []*Payment) {
	payments = []*Payment{}
	if start < 0 || limit < 1 {
		return
	}

	payment_user_id_time_index.lock.RLock()
	defer payment_user_id_time_index.lock.RUnlock()

	iter, _ := payment_user_id_time_index.tree.GetIteratorOrPrev(
		UserIdTimeUnion{userId + 1, 0})
	if iter.Position() != 1 {
		return
	}
	paymentRecord := iter.Value().(PaymentRecord)
	if paymentRecord.GetRef().Payment_user_id == userId+1 {
		if !iter.Prev() {
			return
		}
	}
	if paymentRecord.GetRef().Payment_user_id != userId {
		return
	}
	for ind := 0; ind < start; ind++ {
		if !iter.Prev() {
			return
		}
	}
	sum := 1

	for sum < limit {
		paymentRecord := iter.Value().(PaymentRecord)
		if paymentRecord.GetRef().Payment_user_id != userId {
			return
		}
		payments = append(payments, paymentRecord.GetRef())
		if !iter.Prev() {
			return
		}
		sum++
	}
	return
}

func (item_index *ItemIndex) getItemById(id int) *Item {
	item_index.lock.RLock()
	defer item_index.lock.RUnlock()

	_itemRecord, ok := item_index.tree.Get(id)
	if !ok {
		return nil
	}
	itemRecord, ok := _itemRecord.(ItemRecord)
	if !ok {
		return nil
	}
	return itemRecord.GetRef()
}

func (item_index *ItemIndex) getItemList(start, limit int) (items []*Item) {
	items = []*Item{}
	item_index.lock.RLock()
	defer item_index.lock.RUnlock()

	iter := item_index.tree.Iterator()
	if !iter.Next() {
		return
	}
	for ind := 0; ind < start; ind++ {
		if !iter.Next() {
			return
		}
	}
	for sum := 0; sum < limit; sum++ {
		itemRecord := iter.Value().(ItemRecord)
		items = append(items, itemRecord.GetRef())
		if !iter.Next() {
			return
		}
	}
	return
}

func (uav_index *UAVIndex) getUAVById(id int) *UAV {
	uav_index.lock.RLock()
	defer uav_index.lock.RUnlock()

	_uavRecord, ok := uav_index.tree.Get(id)
	if !ok {
		return nil
	}
	uavRecord, ok := _uavRecord.(UAVRecord)
	if !ok {
		return nil
	}
	return uavRecord.GetRef()
}

func (uav_index *UAVIndex) getUAVList(start, limit int) (uavs []*UAV) {
	uavs = []*UAV{}
	uav_index.lock.RLock()
	defer uav_index.lock.RUnlock()

	iter := uav_index.tree.Iterator()
	if !iter.Next() {
		return
	}
	for ind := 0; ind < start; ind++ {
		if !iter.Next() {
			return
		}
	}
	for sum := 0; sum < limit; sum++ {
		uavRecord := iter.Value().(UAVRecord)
		uavs = append(uavs, uavRecord.GetRef())
		if !iter.Next() {
			return
		}
	}
	return
}

// This function will lock uav for a user
func (uav_index *UAVIndex) getAvailableUAV(userId int) (uav *UAV) {
	uav_index.lock.RLock()
	defer uav_index.lock.RUnlock()

	iter := uav_index.tree.Iterator()
	if !iter.Next() {
		return
	}
	for {
		uavRecord := iter.Value().(UAVRecord)
		uav = uavRecord.GetRef()
		if uav.UAV_status == UAV_STATUS_READY && uav.LockForUser(userId) {
			return uav
		} else {
			if !iter.Next() {
				return nil
			}
		}
	}
}

func ReBuildIndex() {
	user_id_index.lock.Lock()
	user_id_index.tree.Clear()
	for i := 0; i < len(userTable); i++ {
		user_id_index.tree.Put(userTable[i].User_id, UserRecord{
			DB_User: &userTable[i],
		})
	}
	user_id_index.lock.Unlock()

	payment_user_id_time_index.lock.Lock()
	payment_index.lock.Lock()
	payment_user_id_time_index.tree.Clear()
	for i := 0; i < len(paymentTable); i++ {
		key := UserIdTimeUnion{
			paymentTable[i].Payment_user_id, paymentTable[i].Payment_time,
		}
		payment_user_id_time_index.tree.Put(key, PaymentRecord{
			DB_Payment: &paymentTable[i],
		})
		payment_index.tree.Put(paymentTable[i].Payment_id, key)
	}
	payment_index.lock.Unlock()
	payment_user_id_time_index.lock.Unlock()

	item_index.lock.Lock()
	item_index.tree.Clear()
	for i := 0; i < len(itemTable); i++ {
		item_index.tree.Put(itemTable[i].Item_id, ItemRecord{
			DB_Item: &itemTable[i],
		})
	}
	item_index.lock.Unlock()

	uav_index.lock.Lock()
	uav_index.tree.Clear()
	for i := 0; i < len(uavTable); i++ {
		uav_index.tree.Put(uavTable[i].UAV_id, UAVRecord{
			DB_UAV: &uavTable[i],
		})
	}
	uav_index.lock.Unlock()
}
