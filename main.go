package main

import (
	"net/http"

	"github.com/MeteorKL/koala"
)

func main() {
	initFoods()
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	go tcpServer("2018")
	koala.RenderPath = "static/"
	koala.Get("/", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "index.html", nil)
	})
	koala.Get("/user/:id/stop", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "stop_button.html", nil)
	})
	koala.Get("/user/:id/paymentlist", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "paymentlist.html", nil)
	})
	koala.Get("/user/:id/setbutton", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "button.html", nil)
	})
	apiHandlers()
	uavHandlers()
	// http.HandleFunc("/ws", wsHandler)
	koala.Run("2017")
}
