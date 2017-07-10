package main

func initDB() {
	if err := DBLoad("./a.json"); err != nil {
		println(err)
	}
	ReBuildIndex()
}
