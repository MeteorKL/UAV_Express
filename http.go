package main

import (
	"net/http"
	"strconv"

	"github.com/MeteorKL/koala"
)

func apiHandlers() {
	initDB()
	koala.Get("/item", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var start, limit int
		start = koala.GetSingleIntParamOrDefault(p.ParamGet, "start", 0)
		limit = koala.GetSingleIntParamOrDefault(p.ParamGet, "limit", 30)
		items := getItemList(start, limit)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		koala.WriteJSON(w, items)
	})

	koala.Get("/user/:id/payment", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_id := p.ParamUrl["id"]
		num := koala.GetSingleIntParamOrDefault(p.ParamGet, "num", 10)
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		user := getUserById(id)
		if user == nil {
			w.WriteHeader(404)
			w.Write([]byte("no user"))
			return
		}
		payments := user.getRecentPayments(num)
		koala.WriteJSON(w, payments)
	})

	koala.Get("/uav/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		koala.WriteJSON(w, getUAVById(id))
	})

	koala.Get("/uavs", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		uavs := getUAVList(0, 100)
		// var payments []*Payment
		// var users []*User
		var from_to []interface{}
		for id, _ := range uavs {
			if uavs[id].UAV_serving_payment_id != 0 {
				payment := getPaymentById(uavs[id].UAV_serving_payment_id)
				// payments = append(payments, p)
				u := getUserById(payment.Payment_user_id)
				// users = append(users, u)
				from_to = append(from_to, map[string]interface{}{
					"from_longitude": STORE_LONGITUDE,
					"from_latitude":  STORE_LATITUDE,
					"to_longitude":   u.Stop_longitude,
					"to_latitude":    u.Stop_latitude,
				})
			} else {
				from_to = append(from_to, nil)
				// payments = append(payments, nil)
				// users = append(users, nil)
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"status":  0,
			"message": "获取无人机信息成功",
			"data": map[string]interface{}{
				"uavs": uavs,
				// "payments": payments,
				"from_to": from_to,
			},
		},
		)
	})

	koala.Get("/user/:id/payments", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		user := getUserById(id)
		if user == nil {
			w.WriteHeader(404)
			w.Write([]byte("no user"))
			return
		}
		var paymentlist []map[string]interface{}
		var uav interface{}
		payments := user.getRecentPayments(100)
		for _, payment := range payments {
			var itemlist []map[string]interface{}
			var img string
			for id, t := range payment.Payment_items {
				item := getItemById(t.Item_id)
				if id == 0 {
					img = item.Item_img
				}
				itemlist = append(itemlist, map[string]interface{}{
					"name":  item.Item_name,
					"price": item.Item_price * float64(t.Item_num),
					"num":   t.Item_num,
				})
			}
			paymentlist = append(paymentlist, map[string]interface{}{
				"img":   img,
				"id":    payment.Payment_id,
				"time":  payment.Payment_time,
				"items": itemlist,
				"price": payment.Payment_price,
			})
			user := getUserById(payment.Payment_user_id)
			if user == nil {
				w.WriteHeader(400)
				w.Write([]byte("user error"))
				return
			}
			UAV := getUAVById(payment.Payment_uav_id)
			if UAV == nil {
				UAV = nil
			} else {
				uav = map[string]interface{}{
					"status":         UAV.UAV_status,
					"longitude":      UAV.UAV_longitude,
					"latitude":       UAV.UAV_latitude,
					"from_longitude": STORE_LONGITUDE,
					"from_latitude":  STORE_LATITUDE,
					"to_longitude":   user.Stop_longitude,
					"to_latitude":    user.Stop_latitude,
				}
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"uav":         uav,
			"paymentlist": paymentlist,
		})
	})

	koala.Get("/payment/:id/uav", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("param id error"))
			return
		}
		payment := getPaymentById(id)
		if payment == nil {
			w.WriteHeader(400)
			w.Write([]byte("payment error"))
			return
		}
		user := getUserById(payment.Payment_user_id)
		if user == nil {
			w.WriteHeader(400)
			w.Write([]byte("user error"))
			return
		}
		uav := getUAVById(payment.Payment_uav_id)
		if uav == nil {
			w.WriteHeader(400)
			w.Write([]byte("uav error"))
			return
		}
		koala.WriteJSON(w, map[string]interface{}{
			"status":         uav.UAV_status,
			"longitude":      uav.UAV_longitude,
			"latitude":       uav.UAV_latitude,
			"from_longitude": STORE_LONGITUDE,
			"from_latitude":  STORE_LATITUDE,
			"to_longitude":   user.Stop_longitude,
			"to_latitude":    user.Stop_latitude,
		},
		)
	})

	koala.Post("/user/:id/payment", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_items := p.ParamPost["items"]
		_nums := p.ParamPost["nums"]
		if _items == nil || _nums == nil || len(_items) != len(_nums) || len(_items) == 0 {
			w.WriteHeader(400)
			w.Write([]byte("param num error"))
			return
		}
		itemPairs := make([]ItemPair, len(_items))
		var err error
		for i := 0; i < len(_items); i++ {
			itemPairs[i].Item_id, err = strconv.Atoi(_items[i])
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				return
			}
			itemPairs[i].Item_num, err = strconv.Atoi(_nums[i])
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				return
			}
		}

		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		user := getUserById(id)
		if user == nil {
			w.WriteHeader(403)
			w.Write([]byte(err.Error()))
			return
		}

		if user.createPayment(itemPairs) {
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(404)
			w.Write([]byte("No available UAV."))
		}
	})

	koala.Get("/user/:id/button", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("param num error"))
			return
		}
		u := getUserById(id)
		if u == nil {
			w.WriteHeader(400)
			w.Write([]byte("user error"))
			return
		}
		var items []*Item
		for _, id := range u.Stop_buttons {
			item := getItemById(id)
			items = append(items, item)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		koala.WriteJSON(w, items)
	})

	koala.Put("/user/:id/button", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		items := p.ParamPost["items"]
		if items == nil || len(items) != 6 {
			w.WriteHeader(400)
			w.Write([]byte("Not proper number of buttons"))
			return
		}

		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		user := getUserById(id)
		if user == nil {
			w.WriteHeader(404)
			w.Write([]byte("No this user"))
			return
		}
		if user == nil || user.Stop_pin == "" {
			w.WriteHeader(404)
			w.Write([]byte("No this user or user's stop"))
			return
		}
		for i := 0; i < 6; i++ {
			id, err := strconv.Atoi(items[i])
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				return
			}
			user.Stop_buttons[i] = id
		}
		user.Sync()
		w.WriteHeader(200)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("ok"))
	})

	koala.Put("/user/:id/stop", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		pin := koala.GetSingleStringParamOrDefault(p.ParamPost, "pin", "")
		if pin == "" {
			w.WriteHeader(400)
			w.Write([]byte("Not proper string of stop"))
			return
		}

		_id := p.ParamUrl["id"]
		id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		user := getUserById(id)
		if user == nil {
			w.WriteHeader(403)
			w.Write([]byte(err.Error()))
			return
		}

		user.Stop_pin = pin
		user.Sync()
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	koala.Post("/user/:id/stop/pay", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		_userId := p.ParamUrl["id"]
		userId, err := strconv.Atoi(_userId)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("param userId error"))
			return
		}
		user := getUserById(userId)
		if user == nil {
			w.WriteHeader(404)
			w.Write([]byte("No this user"))
			return
		}

		button_id := koala.GetSingleIntParamOrDefault(p.ParamPost, "id", 0)

		if user.createPayment([]ItemPair{{user.Stop_buttons[button_id], 1}}) {
			w.WriteHeader(200)
			koala.WriteJSON(w, map[string]interface{}{
				"message": "ok",
			})
			return
		}
		w.WriteHeader(404)
		w.Write([]byte("No available UAV."))
	})
}
