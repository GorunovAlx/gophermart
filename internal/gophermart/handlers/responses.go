package handlers

type OrderResponse struct {
	Number     string  `json:"number,omitempty"`
	Status     string  `json:"status,omitempty"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at,omitempty"`
}

type BalanceResponse struct {
	Current   float32 `json:"current,omitempty"`
	Withdrawn float32 `json:"withdrawn,omitempty"`
}

type WithdrawResponse struct {
	Order       string  `json:"order,omitempty"`
	Sum         float32 `json:"sum,omitempty"`
	ProcessedAt string  `json:"processed_at"`
}
