package delivery

import (
	"encoding/json"
	"kuncenduit-backend/shared/commandresponse"
	"kuncenduit-backend/shared/httphelper"
	"kuncenduit-backend/shared/lazylogger"
	userSignIn "kuncenduit-backend/src/user/signin/function"
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

	cmd := userSignIn.CommandV1{
		Version:  "1",
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}
	res, errObj := userSignIn.CommandV1Handler(r.Context(), logger, cmd)
	var resObj commandresponse.Obj
	if errObj != nil {
		resObj.Ok = false
		resObj.Err = errObj
	} else {
		jwtCookie := http.Cookie{
			Name:     "jwt",
			Value:    res.SignedJwtToken,
			Expires:  res.JwtExpiration,
			HttpOnly: true, //true will mitigate the risk of client side script accessing the protected cookie
		}
		http.SetCookie(w, &jwtCookie)

		resObj.Ok = true
		resObj.Data = nil
	}

	httphelper.WriteResponseFn(w, resObj)
	httphelper.HandleLogAndPanic(w, logger, errObj)
}
