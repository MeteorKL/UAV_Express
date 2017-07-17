package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MeteorKL/koala"
)

type Specfoods struct {
	Name            string
	Virtual_food_id int
	Price           float64
	Original_price  float64 `json: original_price`
}

type Food struct {
	Name            string
	Description     string
	Image_path      string
	Virtual_food_id int
	Specfoods       []Specfoods
}

type Foods_By_Type struct {
	Type  int
	Name  string
	Foods []Food
}

type Restaurant struct {
	Foods_By_Type []Foods_By_Type
}

//https://www.ele.me/restapi/shopping/v2/menu?restaurant_id=156328250

func CheckErr(info string, err error) {
	if err != nil {
		fmt.Printf("%s: %s", info, err.Error())
	}
}

func WriteToFile(path string, s string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonToFile(path string, jsonStruct interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(jsonStruct); err != nil {
		return err
	}
	return nil
}

// type Test struct {
// 	A string
// 	B string
// }

func initFoods() {
	URL := "https://www.ele.me/restapi/shopping/v2/menu?restaurant_id=156328250"
	statusCode, input := koala.GetRequest(URL)
	if statusCode != 200 {
		fmt.Println("error statusCode: ", statusCode)
		return
	}
	// output, err := iconv.ConvertString(string(input), "gb2312", "utf-8")
	// if err != nil {
	// 	println(URL)
	// 	println("error: iconv.ConvertString()", err.Error())
	// 	// println(string(output))
	// 	return
	// }
	// var t Test
	// err := json.Unmarshal([]byte(`{"a":"1"}`), &t)
	// CheckErr("json error", err)
	// fmt.Println(t)

	var restaurant []Foods_By_Type
	// var output = string(input)

	// WriteToFile("foods", output)
	// output = `{"foods_By_Type":` + output + "}"
	err := json.Unmarshal(input, &restaurant)
	CheckErr("json error", err)

	var items []DB_Item
	for _, Foods_By_Type := range restaurant {
		if Foods_By_Type.Type != 1 {
			continue
		}
		for _, Food := range Foods_By_Type.Foods {
			item := DB_Item{}
			item.Item_name = Food.Name
			item.Item_img = "https://fuss10.elemecdn.com/" + Food.Image_path[0:1] + "/" + Food.Image_path[1:3] + "/" + Food.Image_path[3:] + ".jpeg?imageMogr2/thumbnail/200x200/format/webp/quality/85"

			item.Item_id = Food.Virtual_food_id
			item.Item_description = Food.Description
			item.Item_price = Food.Specfoods[0].Price
			item.Item_type = Foods_By_Type.Name
			items = append(items, item)
		}
	}
	err = WriteJsonToFile("food_result.json", items)
	CheckErr("WriteJsonToFile food_result.json error", err)
}
