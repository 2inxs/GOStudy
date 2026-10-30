package main

import "io"

type Item interface {
	Name() string
	Weight() float64
}

type Usable interface {
	Use() (string, error)
}

// Ошибки чтения и записи передаются вызывающему коду.
type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
