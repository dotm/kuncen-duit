package createerror

import (
	"kuncenduit-backend/shared/commandresponse"
)

func Response(code string, err error, sendErrorTo map[string]bool) *commandresponse.ErrorObj {
	res := &commandresponse.ErrorObj{
		Code:        code,
		SendErrorTo: sendErrorTo,
	}
	if err != nil {
		errMsg := err.Error()
		res.Message = &errMsg
	}
	return res
}

//something that is unexpected and originate from our system
func InternalException(err error) *commandresponse.ErrorObj {
	return Response(
		"internal/exception",
		err,
		map[string]bool{
			commandresponse.SendErrorToLog:      true,
			commandresponse.SendErrorToDevEmail: true,
		},
	)
}

//something that is unexpected and originate from systems other than ours (e.g. partner)
func ExternalException(err error) *commandresponse.ErrorObj {
	return Response(
		"external/exception",
		err,
		map[string]bool{
			commandresponse.SendErrorToLog:      true,
			commandresponse.SendErrorToDevEmail: true,
		},
	)
}

//exceptions caused by bad request from clients.
//if possible don't use this code; instead create new and more domain specific error code
func ClientBadRequest(err error) *commandresponse.ErrorObj {
	return Response(
		"client/bad-request",
		err,
		map[string]bool{
			commandresponse.SendErrorToLog: true,
		},
	)
}
