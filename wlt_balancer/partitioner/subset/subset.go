package subset

type Type struct {
	Numbers []int64
	Sum     int64
}

func New(number int64) *Type {
	return &Type{
		Numbers: []int64{number},
		Sum:     number,
	}
}

func (s *Type) Merge(other *Type) {
	s.Numbers = append(s.Numbers, other.Numbers...)
	s.Sum += other.Sum
}
