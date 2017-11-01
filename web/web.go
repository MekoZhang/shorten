package web

import (
	"log"
	"net/http"

	"github.com/zhangxd1989/shorten/conf"
	"github.com/zhangxd1989/shorten/web/api"
	"github.com/zhangxd1989/shorten/web/www"

	"github.com/gorilla/mux"
)

func Start() {
	log.Println("web starts")
	r := mux.NewRouter()

	r.HandleFunc("/version", api.CheckVersion).Methods(http.MethodGet)
	r.HandleFunc("/health", api.CheckHealth).Methods(http.MethodGet)
	r.HandleFunc("/short", api.ShortURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/expand", api.ExpandURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/{shortenedURL:[a-zA-Z0-9]{1,11}}", api.Redirect).Methods(http.MethodGet)

	r.HandleFunc("/", www.Index).Methods(http.MethodGet)

	r.Handle("/static/{type}/{file}", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(http.Dir("."))))

	if conf.Conf.Common.Schema == "https" {
		log.Fatal(http.ListenAndServeTLS(conf.Conf.Http.Listen, "214310707980181.pem", "214310707980181.key", r));
	} else {
		log.Fatal(http.ListenAndServe(conf.Conf.Http.Listen, r))
	}
}
