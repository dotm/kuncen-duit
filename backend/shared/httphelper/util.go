package httphelper

import (
	"fmt"
	"kuncenduit-backend/shared/commandresponse"
	"kuncenduit-backend/shared/createerror"
	"kuncenduit-backend/shared/lazylogger"
	"net/http"
	"runtime/debug"
)

func HandleLogAndPanic(w http.ResponseWriter, logger *lazylogger.Instance, errObj *commandresponse.ErrorObj) {
	if errObj != nil && errObj.SendErrorTo[commandresponse.SendErrorToLog] {
		logger.LogQueueAsErrorAndDequeueAllItems()
	}
	if err := recover(); err != nil {
		logger.EnqueuePanicLog(err, debug.Stack(), true)
		logger.LogQueueAsErrorAndDequeueAllItems()

		WriteResponseFn(w, commandresponse.Obj{
			Ok:  false,
			Err: createerror.InternalException(fmt.Errorf("pnc")),
		})
	}
}

func WriteResponseFn(w http.ResponseWriter, resObj commandresponse.Obj) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(resObj.ToByteSlice())
}
