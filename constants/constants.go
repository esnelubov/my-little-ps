package constants

import "my-little-ps/helpers/datatypes/set_of_string"

var (
	AllowedCurrencies *set_of_string.Type
)

func Setup() {
	AllowedCurrencies = set_of_string.New("USD", "RUB")
}
