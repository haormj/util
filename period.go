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

// PeriodIntersection intersect two period
func PeriodIntersection(p1 Period, p2 Period) []Period {
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
		return PeriodIntersection(p2, p1)
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

// PeriodDifference b - a
func PeriodDifference(b Period, a Period) []Period {
	switch {
	case a.et < 0 && b.et < 0:
		switch {
		case b.st < a.st:
			ps, _ := NewPeriods(b.st, a.st-1)
			return ps
		case b.st == a.st, b.st > a.st:
			return []Period{}
		}
	case a.et < 0:
		switch {
		case b.et < a.st:
			ps, _ := NewPeriods(b.st, b.et)
			return ps
		case b.et == a.st:
			ps, err := NewPeriods(b.st, b.et-1)
			if err != nil {
				return []Period{}
			}
			return ps
		case b.et > a.st:
			switch {
			case b.st < a.st:
				ps, _ := NewPeriods(b.st, a.st-1)
				return ps
			case b.st == a.st, b.st > a.st:
				return []Period{}
			}
		}
	case b.et < 0:
		switch {
		case b.st < a.st:
			ps, _ := NewPeriods(b.st, a.st-1, a.et+1, b.et)
			return ps
		case b.st == a.st, b.st > a.st && b.st < a.et, b.st == a.et:
			ps, _ := NewPeriods(a.et+1, b.et)
			return ps
		case b.st > a.et:
			ps, _ := NewPeriods(b.st, b.et)
			return ps
		}
	default:
		switch {
		case b.et < a.st:
			ps, _ := NewPeriods(b.st, b.et)
			return ps
		case b.et == a.st:
			ps, err := NewPeriods(b.st, b.et-1)
			if err != nil {
				return []Period{}
			}
			return ps
		case b.et > a.st && b.et < a.et, b.et == a.et:
			ps, err := NewPeriods(b.st, a.st-1)
			if err != nil {
				return []Period{}
			}
			return ps
		case b.et > a.et:
			switch {
			case b.st < a.st:
				ps, _ := NewPeriods(b.st, a.st-1, a.et+1, b.et)
				return ps
			case b.st == a.st, b.st > a.st && b.st < a.et, b.st == a.et:
				ps, _ := NewPeriods(a.et+1, b.et)
				return ps
			case b.st > a.et:
				ps, _ := NewPeriods(b.st, b.et)
				return ps
			}
		}
	}
	return []Period{}
}

// PeriodPartition split period by interval
// period not support et is negative
func PeriodPartition(p Period, interval int64) map[int64]Period {
	m := make(map[int64]Period)
	sti := p.st / interval
	// infinity is not support partition
	if p.et < 0 {
		m[sti*interval] = p
		return m
	}
	eti := p.et / interval
	if sti == eti {
		m[sti*interval] = p
		return m
	}
	stp, _ := NewPeriod(p.st, (sti+1)*interval-1)
	m[sti*interval] = stp

	for i := sti + 1; i < eti; i++ {
		t, _ := NewPeriod(i*interval, (i+1)*interval-1)
		m[i*interval] = t
	}

	etp, _ := NewPeriod(eti*interval, p.et)
	m[eti*interval] = etp

	return m
}

// PeriodMinSuperSet period mininum super set
func PeriodMinSuperSet(p Period, interval int64) Period {
	sti := p.st / interval
	st := sti * interval
	if p.et < 0 {
		ss, _ := NewPeriod(st, p.et)
		return ss
	}
	eti := p.et / interval
	et := (eti+1)*interval - 1
	ss, _ := NewPeriod(st, et)
	return ss
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

// PeriodsIntersection intersect two periods
func PeriodsIntersection(a []Period, b []Period) []Period {
	c := make([]Period, 0)
	for _, i := range a {
		for _, j := range b {
			k := PeriodIntersection(i, j)
			c = append(c, k...)
		}
	}
	return c
}

// PeriodsDifference b - a
func PeriodsDifference(b []Period, a []Period) []Period {
	if len(a) < 1 {
		t := make([]Period, len(b))
		copy(t, b)
		return t
	}

	c := make([]Period, len(b))
	copy(c, b)
	for _, i := range a {
		t := make([]Period, 0)
		for _, j := range c {
			ps := PeriodDifference(j, i)
			t = append(t, ps...)
		}
		c = t
	}
	return c
}

// PeriodsPartition split periods by interval
func PeriodsPartition(ps []Period, interval int64) map[int64][]Period {
	m := make(map[int64][]Period)
	for _, p := range ps {
		t := PeriodPartition(p, interval)
		for k, v := range t {
			mv, ok := m[k]
			if ok {
				mv = append(mv, v)
				m[k] = mv
			} else {
				mv = make([]Period, 0)
				mv = append(mv, v)
				m[k] = mv
			}
		}
	}
	return m
}

// PeriodsMinSuperSet periods mininum super set
func PeriodsMinSuperSet(ps []Period, interval int64) []Period {
	t := make([]Period, 0)
	for _, p := range ps {
		ss := PeriodMinSuperSet(p, interval)
		t = append(t, ss)
	}
	return PeriodsUnion(t, []Period{})
}

// PeriodsSortAsc sort periods by asc
func PeriodsSortAsc(ps []Period) {
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].St() < ps[j].St() {
			return true
		}
		return false
	})
}

// PeriodsSortDesc sort periods by desc
func PeriodsSortDesc(ps []Period) {
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].St() > ps[j].St() {
			return true
		}
		return false
	})
}
