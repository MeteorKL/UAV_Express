package UAV_Express

import (
	"testing"
	"log"
)

func TestDBDump(t *testing.T) {
	itemTable = []DB_Item{
		{1, "咸鱼一条",
			"https://docs.jiguang.cn/jpush/client/image/jpush_android.png",
			100.00, "著名咸鱼S2Meteor的分身",
		},
	}
	userTable = []DB_User{
		{
			1, "S2Meteor",
			"https://docs.jiguang.cn/jpush/client/image/jpush_android.png",
			100.00, "MeteorS2",
			"", 0, 0, [6]int{},
		},
	}
	if err := DBDump("./a.json"); err != nil {
		t.Error(err)
	}
}

func TestDBLoad(t *testing.T) {
	if err := DBLoad("./a.json"); err != nil {
		t.Error(err)
	}
	log.Println(userTable)
}