package main

import (
	"fmt"
	userSigninDelivery "kuncenduit-backend/src/user/signin/delivery"
	"kuncenduit-e2e-backend/shared"
	"net/http"
)

func main() {
	endpoint := shared.BackendUrl + userSigninDelivery.PathV1
	body := userSigninDelivery.DTOV1{
		Email:    "local-kd@yopmail.com",
		Password: "Test123!",
	}
	req, err := shared.CreatePostRequestWithJsonBody(body, endpoint)
	if err != nil {
		return //error already logged from shared util function
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error executing request: %v", err.Error())
		return
	}
	defer res.Body.Close()

	err = shared.LogResponse(res)
	if err != nil {
		return //error already logged from shared util function
	}
}
