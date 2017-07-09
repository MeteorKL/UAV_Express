package main

import (
	"testing"

	"github.com/MeteorKL/koala"
)

// go test -v http_test.go api.go db.go index.go

func tInitItemTable() {
	itemTable = []DB_Item{
		{1, "咸鱼一条",
			"https://docs.jiguang.cn/jpush/client/image/jpush_android.png",
			100.00, "著名咸鱼S2Meteor的分身",
		},
	}
	ReBuildIndex()
}

func tGetItemList() string {
	return string(koala.GetRequest("http://localhost:2017/api/item"))
}

func Test_getItemList(t *testing.T) {
	tInitItemTable()
	t.Log(tGetItemList())
}
