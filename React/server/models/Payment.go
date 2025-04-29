package models

type Payment struct {
	Amount           float64 `json:"amount"`
	Name             string  `json:"name"`
	Surname          string  `json:"surname"`
	CreditCardNumber string  `json:"credit_card_number"`
}
