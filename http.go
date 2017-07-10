package main

import (
	"net/http"
	"strconv"

	"github.com/MeteorKL/koala"
)

func apiHandlers() {
	initDB()
	koala.Get("/api/item", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var start, limit int
		start = koala.GetSingleIntParamOrDefault(p.ParamGet, "start", 0)
		limit = koala.GetSingleIntParamOrDefault(p.ParamGet, "limit", 30)
		items := getItemList(start, limit)
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
			w.WriteHeader(403)
			w.Write([]byte(err.Error()))
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

	koala.Post("/stop/:pin/pay/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		pin := p.ParamUrl["pin"]
		userId := koala.GetSingleIntParamOrDefault(p.ParamPost, "userId", 0)
		user := getUserById(userId)
		if user == nil || user.Stop_pin == pin {
			w.WriteHeader(404)
			w.Write([]byte("No this user or pin or not matching"))
			return
		}

		_id := p.ParamUrl["id"]
		button_id, err := strconv.Atoi(_id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		if user.createPayment([]ItemPair{{user.Stop_buttons[button_id], 1}}) {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
			return
		}
		w.WriteHeader(404)
		w.Write([]byte("No available UAV."))
	})
}
