package set_of_string

import "fmt"

type Type map[string]bool

func (s *Type) Items() []string {
	keys := make([]string, len(*s))

	i := 0
	for k := range *s {
		keys[i] = k
		i++
	}

	return keys
}

func New(values ...string) (result *Type) {
	result = &Type{}

	for _, val := range values {
		result.Add(val)
	}

	return
}

func (s *Type) Add(value string) {
	(*s)[value] = true
}

func (s *Type) Has(value string) bool {
	_, ok := (*s)[value]

	return ok
}

func (s *Type) Empty() bool {
	return len(*s) == 0
}

func (s *Type) String() string {
	return fmt.Sprintf("%v", s.Items())
}
