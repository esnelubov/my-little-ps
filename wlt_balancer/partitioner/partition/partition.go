package partition

import (
	"my-little-ps/wlt_balancer/partitioner/subset"
	"sort"
)

type Type struct {
	Subsets []*subset.Type
}

func New(number int64, k int) *Type {
	p := &Type{
		Subsets: make([]*subset.Type, k),
	}

	p.Subsets[0] = subset.New(number)
	return p
}

func (p *Type) Priority() int64 {
	first := p.Subsets[0]
	last := p.Subsets[len(p.Subsets)-1]

	if first == nil {
		return 0
	}

	if last == nil {
		return first.Sum
	}

	return first.Sum - last.Sum
}

func (p *Type) Merge(other *Type) {
	for i, a := range p.Subsets {
		b := other.Subsets[len(other.Subsets)-1-i]

		if a == nil {
			p.Subsets[i] = b
			continue
		}

		if b == nil {
			continue
		}

		a.Merge(b)
	}

	sort.Slice(p.Subsets, func(i, j int) bool {
		a := p.Subsets[i]
		b := p.Subsets[j]

		if a == nil && b == nil {
			return false
		}

		if a == nil && b != nil {
			return false
		}

		if a != nil && b == nil {
			return true
		}

		return a.Sum > b.Sum
	})
}
