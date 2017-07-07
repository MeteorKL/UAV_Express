package main

import (
	"net/http"

	"github.com/MeteorKL/koala"
)

func main() {
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))

	koala.RenderPath = "static/"
	koala.Get("/", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "index.html", nil)
	})
	uavHandlers()
	// http.HandleFunc("/ws", wsHandler)
	koala.RunWithLog("2017")
}
