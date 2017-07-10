package main

import (
	"errors"
	"strconv"
	"testing"

	"encoding/json"

	"github.com/MeteorKL/koala"
)

// go test -v http_test.go api.go db.go index.go

func tGetItemList() string {
	_, data := koala.Request("GET", "http://localhost:2017/api/item", "")
	return string(data)
}

func tGetPayments(id string) (payments []DB_Payment, err error) {
	status, data := koala.Request("GET", "http://localhost:2017/user/"+id+"/payment", "")
	if status == 200 {
		paymentData := make(map[string][]DB_Payment)
		err = json.Unmarshal([]byte("{\"data\":"+string(data)+"}"), &paymentData)
		payments = paymentData["data"]
	} else {
		err = errors.New(strconv.Itoa(status))
	}
	return
}

func tPostPayment(id string) string {
	_, data := koala.Request("POST", "http://localhost:2017/user/"+id+"/payment", "items=1&nums=1")
	return string(data)
}

func tGetUav(id string) string {
	_, data := koala.Request("GET", "http://localhost:2017/uav/"+id, "")
	return string(data)
}
func tGetUavs() string {
	_, data := koala.Request("GET", "http://localhost:2017/uavs", "")
	return string(data)
}

func tGetPayment(id string) {

}

func Test_getItemList(t *testing.T) {
	t.Log(tGetItemList())
	if payments, err := tGetPayments("1"); err == nil {
		t.Log("订单数量: ", len(payments))
	}
	t.Log(tPostPayment("1"))
	if payments, err := tGetPayments("1"); err == nil && len(payments) > 0 {
		t.Log("订单数量: ", len(payments))
		for i := range payments {
			t.Log(payments[i])
			t.Log(getPaymentById(payments[i].Payment_id))
			_uav_id := payments[i].Payment_uav_id
			uav_id := strconv.Itoa(_uav_id)
			t.Log(tGetUav(uav_id))
		}
	}
	t.Log(tGetUavs())
}
