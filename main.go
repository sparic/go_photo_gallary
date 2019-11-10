package main

import (
	"scratch_maker_server/routers"
	"net/http"
	"time"
)

func main() {
	// get the global router from router.go
	router := routers.Router
	// router := gin.Default()

	// set up a http server
	server := http.Server{
		// Addr:           fmt.Sprintf(":%s", constant.SERVER_PORT),
		Addr:           ":1234",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// run the server
	server.ListenAndServe()
}
