package delivery

const PathV1 = "/v1/user.signup"

type DTOV1 struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
