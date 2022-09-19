package delivery

import (
	"encoding/json"
	"kuncenduit-backend/shared/commandresponse"
	"kuncenduit-backend/shared/httphelper"
	"kuncenduit-backend/shared/lazylogger"
	userSignUp "kuncenduit-backend/src/user/signup/function"
	"net/http"
)

func HttpHandlerV1(w http.ResponseWriter, r *http.Request) {
	//logging and panic handling needs to be copied
	//for all delivery methods (HTTP server, Serverless Function, etc.)
	var errObj *commandresponse.ErrorObj
	logger := lazylogger.New(r.URL.Path)

	var reqBody DTOV1

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := userSignUp.CommandV1{
		Version:  "1",
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}
	res, errObj := userSignUp.CommandV1Handler(r.Context(), logger, cmd)
	var resObj commandresponse.Obj
	if errObj != nil {
		resObj.Ok = false
		resObj.Err = errObj
	} else {
		resObj.Ok = true
		resObj.Data = res
	}

	httphelper.WriteResponseFn(w, resObj)
	httphelper.HandleLogAndPanic(w, logger, errObj)
}
