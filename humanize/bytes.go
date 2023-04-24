package humanize

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// IEC Sizes.
// kibis of bits
const (
	Byte = 1 << (iota * 10)
	KiByte
	MiByte
	GiByte
	TiByte
	PiByte
	EiByte
)

type Bytes int64

var bytesRE = regexp.MustCompile("^(([0-9]+)eib)?(([0-9]+)pib)?(([0-9]+)tib)?(([0-9]+)gib)?(([0-9]+)mib)?(([0-9]+)kib)?(([0-9]+)b)?$")

// ParseBytes parses a string into a Bytes
func ParseBytes(bytesStr string) (Bytes, error) {
	switch bytesStr {
	case "", "0":
		// Allow 0 without a unit.
		return 0, nil
	}
	matches := bytesRE.FindStringSubmatch(bytesStr)
	if matches == nil {
		return 0, fmt.Errorf("not a valid bytes string: %q", bytesStr)
	}
	var size Bytes

	// Parse the match at pos `pos` in the regex and use `mult` to turn that
	// into ms, then add that value to the total parsed duration.
	var overflowErr error
	m := func(pos int, mult Bytes) {
		if matches[pos] == "" {
			return
		}
		n, _ := strconv.Atoi(matches[pos])

		size += Bytes(n) * mult

		if size < 0 {
			overflowErr = errors.New("bytes out of range")
		}
	}

	m(2, EiByte)  // EiB
	m(4, PiByte)  // PiB
	m(6, TiByte)  // TiB
	m(8, GiByte)  // GiB
	m(10, MiByte) // MiB
	m(12, KiByte) // KiB
	m(14, Byte)   // B

	return Bytes(size), overflowErr
}

func (b Bytes) String() string {
	var (
		n = int64(b)
		r = ""
	)
	if n == 0 {
		return "0b"
	}

	f := func(unit string, mult int64, exact bool) {
		if exact && n%mult != 0 {
			return
		}
		if v := n / mult; v > 0 {
			r += fmt.Sprintf("%d%s", v, unit)
			n -= v * mult
		}
	}

	f("eib", EiByte, false)
	f("pib", PiByte, false)
	f("tib", TiByte, false)
	f("gib", GiByte, false)
	f("mib", MiByte, false)
	f("kib", KiByte, false)
	f("b", Byte, false)

	return r
}

// MarshalJSON implements the json.Marshaler interface.
func (b Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bytes) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	bs, err := ParseBytes(s)
	if err != nil {
		return err
	}
	*b = bs
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (b *Bytes) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *Bytes) UnmarshalText(text []byte) error {
	var err error
	*b, err = ParseBytes(string(text))
	return err
}

// MarshalYAML implements the yaml.Marshaler interface.
func (b Bytes) MarshalYAML() (interface{}, error) {
	return b.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (b *Bytes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	bs, err := ParseBytes(s)
	if err != nil {
		return err
	}
	*b = bs
	return nil
}
