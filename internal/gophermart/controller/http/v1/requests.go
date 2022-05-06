package v1

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type WithdrawRequest struct {
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
}
