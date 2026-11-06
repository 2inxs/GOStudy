package items

import "io"

type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
}

type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
