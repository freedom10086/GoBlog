package main

import (
	"goBlog/conf"
	"goBlog/model"
	"goBlog/router"
	"log"
	"net/http"
	"time"
)

//配置文件
var config *conf.Config

//路由列表
var routers map[string]router.Handler

func init() {
	config = conf.Conf
	model.InitDB(config.DbName, config.DbUsername, config.DbPassword)
	log.Printf("==%s started==", config.SiteName)

	routers = map[string]router.Handler{
		"/":          &router.HomeHandler{},
		"/static/":   &router.StaticFileHandler{},
		"/categorys": &router.CateHandler{},
		"/users":     &router.UserHandler{},
		"/auth":      &router.OauthHandler{},
		"/login":     &router.LoginHandler{},
		"/register":  &router.RegisterHandler{},
		"/chats":     &router.ChatHandler{},
	}
}

func main() {
	defer model.CloseDB()
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := "https://127.0.0.1" + config.SitePortSSL + r.URL.Path
			router.Redirect(w, r, path, http.StatusMovedPermanently)
		})
		http.ListenAndServe(config.SitePort, nil)
	}()

	r := router.NewRouter()
	for key, value := range routers {
		r.Register(key, value)
	}

	server := &http.Server{
		Addr:           config.SitePortSSL,
		Handler:        r,
		WriteTimeout:   8 * time.Second,
		ReadTimeout:    8 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("https listen on %s%s", "https://127.0.0.1", config.SitePortSSL)
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
