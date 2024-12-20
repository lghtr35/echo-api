package operation

type Operation uint

type Operable[T any] struct {
	Op  Operation
	Val T
}

const (
	Insert Operation = iota
	Delete
)
