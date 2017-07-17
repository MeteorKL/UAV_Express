package main

import (
	"bufio"
	"encoding/json"
	"os"
)

var (
	itemTable    = []DB_Item{}
	userTable    = []DB_User{}
	uavTable     = []DB_UAV{}
	paymentTable = []DB_Payment{}
)

const (
	STORE_LONGITUDE = 120.131658
	STORE_LATITUDE  = 30.273568
)

type DB_Item struct {
	Item_id          int     `json:"item_id"`
	Item_name        string  `json:"item_name"`
	Item_type        string  `json:"item_type"`
	Item_img         string  `json:"item_img"`
	Item_price       float64 `json:"item_price"`
	Item_description string  `json:"item_description"`
}

type DB_User struct {
	User_id        int     `json:"user_id"`
	User_name      string  `json:"user_name"`
	User_avatar    string  `json:"user_avatar"`
	User_balance   float64 `json:"user_balance"`
	User_password  string  `json:"user_password"`
	Stop_pin       string  `json:"stop_pin"`
	Stop_longitude float64 `json:"stop_longitude"`
	Stop_latitude  float64 `json:"stop_latitude"`
	Stop_buttons   [6]int  `json:"stop_buttons"`
}

const (
	UAV_STATUS_READY = iota
	UAV_STATUS_SENDING
	UAV_STATUS_LANDING
	UAV_STATUS_RETURNING
)

type DB_UAV struct {
	UAV_id                 int     `json:"uav_id"`
	UAV_name               string  `json:"uav_name"`
	UAV_longitude          float64 `json:"uav_longitude"`
	UAV_latitude           float64 `json:"uav_latitude"`
	UAV_serving_payment_id int     `json:"uav_serving_payment_id"`
	UAV_status             int32   `json:"uav_status"`
}

type ItemPair struct {
	Item_id  int `json:"item_id"`
	Item_num int `json:"item_num"`
}

type DB_Payment struct {
	Payment_id      int        `json:"payment_id"`
	Payment_number  string     `json:"payment_number"`
	Payment_items   []ItemPair `json:"payment_item_id"`
	Payment_user_id int        `json:"payment_user_id"`
	Payment_uav_id  int        `json:"payment_uav_id"`
	Payment_price   float64    `json:"payment_price"`
	Payment_time    int        `json:"payment_time"`
}

func DBDump(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(itemTable); err != nil {
		return err
	}
	if err := encoder.Encode(userTable); err != nil {
		return err
	}
	if err := encoder.Encode(uavTable); err != nil {
		return err
	}
	if err := encoder.Encode(paymentTable); err != nil {
		return err
	}
	return nil
}

func DBLoad(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&userTable); err != nil {
		return err
	}
	if err := decoder.Decode(&uavTable); err != nil {
		return err
	}
	if err := decoder.Decode(&paymentTable); err != nil {
		return err
	}
	return nil
}

func DBLoadItem(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&itemTable); err != nil {
		return err
	}
	return nil
}
