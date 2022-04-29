package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/handlers"
)

func main() {
	config.SetConfig()
	router := handlers.Initialize()
	srv := &http.Server{
		Handler:      router.Negroni,
		Addr:         config.Cfg.RunAddress,
		WriteTimeout: 1015000 * time.Second,
		ReadTimeout:  1015000 * time.Second,
		//IdleTimeout:  time.Second * 60 * 5,
	}

	log.Fatal(srv.ListenAndServe())
}
