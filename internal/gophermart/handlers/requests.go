package handlers

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BalanceRequest struct {
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
}
