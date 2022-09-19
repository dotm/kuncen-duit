package delivery

const PathV1 = "/v1/user.signin"

type DTOV1 struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
