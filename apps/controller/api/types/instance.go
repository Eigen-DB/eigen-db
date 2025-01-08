package types

type CreateBody struct {
	CustomerId string `json:"customerId" validate:"required"`
	// add requested resources...
}

type TerminateBody struct {
	CustomerId string `json:"customerId" validate:"required"`
}
