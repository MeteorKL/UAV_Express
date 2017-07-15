package main

func initDB() {
	if err := DBLoad("./a.json"); err != nil {
		println(err)
	}
	if err := DBLoadItem("./item.json"); err != nil {
		println(err)
	}
	ReBuildIndex()
}
