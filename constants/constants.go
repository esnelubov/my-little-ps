package constants

import "my-little-ps/helpers/datatypes/set_of_string"

var (
	AllowedCurrencies *set_of_string.Type
)

const (
	OpStatusNew        = "new"
	OpStatusProcessing = "processing"
	OpStatusSuccess    = "success"
	OpStatusDecline    = "decline"
)

func Setup() {
	AllowedCurrencies = set_of_string.New("USD", "RUB")
}
