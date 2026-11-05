package items

import "io"

type Item interface {
	ID() int
	Name() string
	Weight() float64
}

type Usable interface {
	Use() (string, error)
}

type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
