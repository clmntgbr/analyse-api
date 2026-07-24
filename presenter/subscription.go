package presenter

type CheckoutSessionResponse struct {
	URL string `json:"url"`
}

func NewCheckoutSessionResponse(url string) CheckoutSessionResponse {
	return CheckoutSessionResponse{URL: url}
}
