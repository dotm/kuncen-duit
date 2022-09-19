package main

import (
	"fmt"
	"io"
	userSigninDelivery "kuncenduit-backend/src/user/signin/delivery"
	userSignupDelivery "kuncenduit-backend/src/user/signup/delivery"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

var port = ":8080"

func main() {
	log.Println("run main")

	loadEnv()

	http.HandleFunc("/ping", bypassCORS(ping))
	http.HandleFunc(userSignupDelivery.PathV1, bypassCORS(userSignupDelivery.HttpHandlerV1))
	http.HandleFunc(userSigninDelivery.PathV1, bypassCORS(userSigninDelivery.HttpHandlerV1))
	log.Println("listening on", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Print(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "pong\n")
}

func loadEnv() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}
}
