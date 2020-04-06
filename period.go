package util

import (
	"errors"
	"sort"
)

var (
	ErrPeriodInvalid    = errors.New("invalid period")
	ErrPeriodStNegative = errors.New("start time is negative")
	ErrPeriodIllegalArg = errors.New("illegal argument")
)

// Period a length or portion of time.
// [StartTime, EndTime]
// StartTime must be >= 0
// EndTime negative values represent positive infinity
type Period struct {
	st int64
	et int64
}

func (p Period) St() int64 {
	return p.st
}

func (p Period) Et() int64 {
	return p.et
}

type Periods []Period

// Len is the number of elements in the collection.
func (p Periods) Len() int {
	return len(p)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (p Periods) Less(i, j int) bool {
	if p[i].st < p[j].st {
		return true
	}
	return false
}

// Swap swaps the elements with indexes i and j.
func (p Periods) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// NewPeriod create a period
func NewPeriod(st, et int64) (Period, error) {
	if st < 0 {
		return Period{}, ErrPeriodStNegative
	}
	if et >= 0 && st > et {
		return Period{}, ErrPeriodInvalid
	}
	p := Period{
		st: st,
		et: et,
	}
	return p, nil
}

// NewPeriods create a set of period
func NewPeriods(se ...int64) ([]Period, error) {
	if len(se)%2 != 0 {
		return nil, ErrPeriodIllegalArg
	}
	r := make([]Period, 0)
	for i := 0; i < len(se); i = i + 2 {
		p, err := NewPeriod(se[i], se[i+1])
		if err != nil {
			return nil, err
		}
		r = append(r, p)
	}
	return r, nil
}

// PeriodContains judge is t in period
func PeriodContains(t int64, p Period) bool {
	if t < 0 {
		return false
	}
	if t >= p.st {
		if p.et < 0 {
			return true
		}
		if t <= p.et {
			return true
		}
	}
	return false
}

// PeriodUnion union two period
func PeriodUnion(p1 Period, p2 Period) []Period {
	switch {
	case p1.et < 0 && p2.et < 0:
		st := p2.st
		et := p2.et
		if p1.st < p2.st {
			st = p1.st
		}
		ps, _ := NewPeriods(st, et)
		return ps
	case p1.et < 0:
		switch {
		case p2.et < p1.st:
			return []Period{p1, p2}
		case p2.et == p1.st:
			st := p2.st
			et := p1.et
			ps, _ := NewPeriods(st, et)
			return ps
		case p2.et > p1.st:
			st := p1.st
			et := p1.et
			if p2.st < p1.st {
				st = p2.st
			}
			ps, _ := NewPeriods(st, et)
			return ps
		}
	case p2.et < 0:
		return PeriodUnion(p2, p1)
	default:
		switch {
		case p2.et < p1.st:
			return []Period{p1, p2}
		case p2.et == p1.st:
			st := p2.st
			et := p1.et
			ps, _ := NewPeriods(st, et)
			return ps
		case p2.et > p1.st && p2.et < p1.et, p2.et == p1.et:
			st := p1.st
			et := p1.et
			if p2.st < p1.st {
				st = p2.st
			}
			ps, _ := NewPeriods(st, et)
			return ps
		case p2.et > p1.et:
			switch {
			case p2.st < p1.st, p2.st == p1.st:
				return []Period{p2}
			case p2.st > p1.st && p2.st < p1.et, p2.st == p1.et:
				st := p1.st
				et := p2.et
				ps, _ := NewPeriods(st, et)
				return ps
			case p2.st > p1.et:
				return []Period{p1, p2}
			}
		}
	}
	return []Period{p1, p2}
}

// PeriodIntersect intersect two period
func PeriodIntersect(p1 Period, p2 Period) []Period {
	switch {
	case p2.et < 0 && p1.et < 0:
		st := p2.st
		et := p2.et
		if p2.st < p1.st {
			st = p1.st
		}
		ps, _ := NewPeriods(st, et)
		return ps
	case p1.et < 0:
		switch {
		case p2.et < p1.st:
			return []Period{}
		case p2.et == p1.st:
			st := p2.et
			et := p2.et
			ps, _ := NewPeriods(st, et)
			return ps
		case p2.et > p1.st:
			st := p2.st
			et := p2.et
			if p2.st < p1.st {
				st = p1.st
			}
			ps, _ := NewPeriods(st, et)
			return ps
		}
	case p2.et < 0:
		return PeriodIntersect(p2, p1)
	default:
		switch {
		case p2.et < p1.st:
			return []Period{}
		case p2.et == p1.st:
			st := p2.et
			et := p2.et
			ps, _ := NewPeriods(st, et)
			return ps
		case p2.et > p1.st && p2.et < p1.et, p2.et == p1.et:
			st := p2.st
			et := p2.et
			if p2.st < p1.st {
				st = p1.st
			}
			ps, _ := NewPeriods(st, et)
			return ps
		case p2.et > p1.et:
			switch {
			case p2.st < p1.st, p2.st == p1.st:
				st := p1.st
				et := p1.et
				ps, _ := NewPeriods(st, et)
				return ps
			case p2.st > p1.st && p2.st < p1.et:
				st := p2.st
				et := p1.et
				ps, _ := NewPeriods(st, et)
				return ps
			case p2.st == p1.et:
				st := p2.st
				et := p2.st
				ps, _ := NewPeriods(st, et)
				return ps
			case p2.st > p1.et:
				return []Period{}
			}
		}
	}
	return []Period{}
}

// PeriodsContains judge is t in periods
func PeriodsContains(t int64, ps []Period) bool {
	for _, p := range ps {
		if PeriodContains(t, p) {
			return true
		}
	}
	return false
}

// AddPeriodToResultSet add a period to result set
// ResultSet is a set of period, all elements has been union
func AddPeriodToResultSet(p Period, rs []Period) []Period {
	if len(rs) == 0 {
		return []Period{p}
	}
	for i := 0; i < len(rs); i++ {
		ps := PeriodUnion(p, rs[i])
		if len(ps) == 1 {
			t := rs[:i]
			if i+1 < len(rs) {
				t = append(t, rs[i+1:]...)
			}
			return AddPeriodToResultSet(ps[0], t)
		}
	}
	return append(rs, p)
}

// PeriodsUnion union two periods
func PeriodsUnion(a []Period, b []Period) []Period {
	ps := append(a, b...)
	rs := make([]Period, 0)
	for _, p := range ps {
		rs = AddPeriodToResultSet(p, rs)
	}
	return rs
}

// PeriodsIntersect intersect two periods
func PeriodsIntersect(a []Period, b []Period) []Period {
	c := make([]Period, 0)
	for _, i := range a {
		for _, j := range b {
			k := PeriodIntersect(i, j)
			c = append(c, k...)
		}
	}
	return c
}

// PeriodsComplement periods complement
// u = [a.min_st, a.max_et]
// https://en.wikipedia.org/wiki/Complement_(set_theory)
func PeriodsComplement(a []Period) []Period {
	if len(a) < 1 {
		return []Period{}
	}
	b := PeriodsUnion(a, []Period{})
	sort.Sort(Periods(b))
	c := make([]Period, 0)
	for i := 0; i < len(b)-1; i++ {
		st := b[i].et + 1
		// et not negative
		et := b[i+1].st - 1
		p, err := NewPeriod(st, et)
		if err != nil {
			continue
		}
		c = append(c, p)
	}
	return c
}
