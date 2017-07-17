package main

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

// These four structures is copy of DB-data
type User struct {
	DB_User
}

type Payment struct {
	DB_Payment
}

type Item struct {
	DB_Item
}

type UAV struct {
	DB_UAV
	lock_user_id chan int `json:"-"`
}

const v float64 = 0.030

func (uav *UAV) move(user_id int, from_longitude float64, from_latitude float64, to_longitude float64, to_latitude float64) {
	distance_from_to := distance(from_longitude, from_latitude, to_longitude, to_latitude)
	r := v / distance_from_to
	u := getUserById(user_id)
	if u.User_balance < 0 {
		sendMsg2User(user_id, "您已欠费"+strconv.FormatFloat(-u.User_balance, 'f', 2, 64)+"，请及时充值")
	}
	sendMsg2User(user_id, "无人机已经出发")
	for {
		timer1 := time.NewTicker(time.Second)
		for {
			select {
			case <-timer1.C:
				switch uav.UAV_status {
				case UAV_STATUS_SENDING:
					uav.UAV_longitude = r*(to_longitude-from_longitude) + uav.UAV_longitude
					uav.UAV_latitude = r*(to_latitude-from_latitude) + uav.UAV_latitude
					if distance(uav.UAV_longitude, uav.UAV_latitude, from_longitude, from_latitude) > distance_from_to {
						uav.UAV_longitude = to_longitude
						uav.UAV_latitude = to_latitude
						uav.UAV_status = UAV_STATUS_LANDING
						sendMsg2User(user_id, "无人机已经到达，正在着陆")
						println("reached")
					}
					uav.Sync()
				case UAV_STATUS_LANDING:
					payment := getPaymentById(uav.UAV_serving_payment_id)
					price := payment.Payment_price
					user := getUserById(payment.Payment_user_id)
					user.User_balance -= price
					user.Sync()
					sendMsg2User(user_id, "订单已完成，本次消费"+strconv.FormatFloat(price, 'f', 2, 64)+"，账户余额"+strconv.FormatFloat(user.User_balance, 'f', 2, 64))
					time.Sleep(time.Second * 5)
					uav.UAV_status = UAV_STATUS_RETURNING
					uav.Sync()
				case UAV_STATUS_RETURNING:
					uav.UAV_longitude = r*(from_longitude-to_longitude) + uav.UAV_longitude
					uav.UAV_latitude = r*(from_latitude-to_latitude) + uav.UAV_latitude
					if distance(uav.UAV_longitude, uav.UAV_latitude, to_longitude, to_latitude) > distance_from_to {
						uav.UAV_longitude = from_longitude
						uav.UAV_latitude = from_latitude
						uav.UAV_status = UAV_STATUS_READY
						uav.UAV_serving_payment_id = 0
						println("finished")
					}
					uav.Sync()
				case UAV_STATUS_READY:
					return
				}
			}
		}
	}
}

func (uav *UAV) LockForUser(userId int) bool {
	select {
	case uav.lock_user_id <- userId:
		return true
	default:
		return false
	}
}

func (uav *UAV) UnLock() bool {
	select {
	case <-uav.lock_user_id:
		return true
	default:
		return false
	}
}

func getUserById(id int) *User {
	return user_id_index.getUserById(id)
}

func (user *User) getRecentPayments(num int) []*Payment {
	return payment_user_id_time_index.getUserLastPayments(user.User_id, 0, num)
}

func getItemById(id int) *Item {
	return item_index.getItemById(id)
}

func getItemList(start, limit int) []*Item {
	return item_index.getItemList(start, limit)
}

func getUAVById(id int) *UAV {
	return uav_index.getUAVById(id)
}

func getUAVList(start, limit int) []*UAV {
	return uav_index.getUAVList(start, limit)
}

func (user *User) getAvailableUAV() *UAV {
	return uav_index.getAvailableUAV(user.User_id)
}

func getPaymentById(id int) *Payment {
	return payment_index.getPayment(id)
}

func (user *User) createPayment(pairs []ItemPair) bool {
	price := 0.0
	for _, p := range pairs {
		item := getItemById(p.Item_id)
		price += item.Item_price * float64(p.Item_num)
	}

	uav := uav_index.getAvailableUAV(user.User_id)
	if uav == nil {
		return false
	}
	uav.UAV_status = UAV_STATUS_SENDING
	Payment_id := int(rand.Int31()) //TODO: Sync problem?
	uav.UAV_serving_payment_id = Payment_id
	uav.Sync()
	go uav.move(user.User_id, uav.UAV_longitude, uav.UAV_latitude, user.Stop_longitude, user.Stop_latitude)
	defer uav.UnLock()

	payment := Payment{
		DB_Payment{
			Payment_id:      Payment_id,
			Payment_time:    int(time.Now().UnixNano() / 1000000),
			Payment_price:   price,
			Payment_user_id: user.User_id,
			Payment_items:   pairs,
			Payment_number:  "S2Meteor", //TODO:
			Payment_uav_id:  uav.UAV_id,
		},
	}
	payment_user_id_time_index.insertPayment(payment)
	if payment.Sync() != nil {
		return false
	}
	return true
}

func (user *User) Sync() error {
	user_id_index.lock.Lock()
	defer user_id_index.lock.Unlock()
	_userRecord, ok := user_id_index.tree.Get(user.User_id)
	if !ok {
		return errors.New("This user has been deleted!")
	}
	userRecord := _userRecord.(UserRecord)
	*userRecord.DB_User = user.DB_User
	return nil
}

func (payment *Payment) Sync() error {
	payment_user_id_time_index.lock.Lock()
	defer payment_user_id_time_index.lock.Unlock()
	payment_index.lock.Lock()
	defer payment_index.lock.Unlock()

	_paymentRecord, ok := payment_user_id_time_index.tree.Get(UserIdTimeUnion{
		payment.Payment_user_id, payment.Payment_time,
	})
	if !ok {
		return errors.New("This payment has been deleted!")
	}
	paymentRecord := _paymentRecord.(PaymentRecord)
	*paymentRecord.DB_Payment = payment.DB_Payment
	return nil
}

func (item *Item) Sync() error {
	item_index.lock.Lock()
	defer item_index.lock.Unlock()
	_itemRecord, ok := item_index.tree.Get(item.Item_id)
	if !ok {
		return errors.New("This item has been deleted!")
	}
	itemRecord := _itemRecord.(ItemRecord)
	*itemRecord.DB_Item = item.DB_Item
	return nil
}

func (uav *UAV) Sync() error {
	uav_index.lock.Lock()
	defer uav_index.lock.Unlock()
	_uavRecord, ok := uav_index.tree.Get(uav.UAV_id)
	if !ok {
		return errors.New("This uav has been deleted!")
	}
	uavRecord := _uavRecord.(UAVRecord)
	*uavRecord.DB_UAV = uav.DB_UAV
	return nil
}
