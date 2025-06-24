package safecast

// This files is highly inspired from https://pkg.go.dev/golang.org/x/exp/constraints
// I didn't import it as it would have added an unnecessary dependency

// Signed is a constraint for all signed integers: int, int8, int16, int32, and int64 types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint for all unsigned integers: uint, uint8, uint16, uint32, and uint64 types.
// TODO: support uintpr
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Integer is a constraint for the all unsigned and signed integers
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint for the float32 and float64 types.
type Float interface {
	~float32 | ~float64
}

// Number is a constraint for all integers and floats
// TODO: consider using complex, but not sure there is a need
type Number interface {
	Integer | Float
}

// Input is a constraint for all types that can be used as input for [Convert], and [MustConvert]
type Input interface {
	Number | ~string | ~bool
}
