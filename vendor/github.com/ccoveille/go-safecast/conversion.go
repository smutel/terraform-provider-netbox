package safecast

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Convert attempts to convert any value to the desired type
//   - If the conversion is possible, the converted value is returned.
//   - If the conversion results in a value outside the range of the desired type, an [ErrRangeOverflow] error is wrapped in the returned error.
//   - If the conversion exceeds the maximum value of the desired type, an [ErrExceedMaximumValue] error is wrapped in the returned error.
//   - If the conversion exceeds the minimum value of the desired type, an [ErrExceedMinimumValue] error is wrapped in the returned error.
//   - If the conversion is not possible for the desired type, an [ErrUnsupportedConversion] error is wrapped in the returned error.
//   - If the conversion fails from string, an [ErrStringConversion] error is wrapped in the returned error.
//   - If the conversion results in an error, an [ErrConversionIssue] error is wrapped in the returned error.
func Convert[NumOut Number, NumIn Input](orig NumIn) (converted NumOut, err error) {
	v := reflect.ValueOf(orig)
	switch v.Kind() {
	case reflect.Int:
		return convertFromNumber[NumOut](int(v.Int()))
	case reflect.Uint:
		return convertFromNumber[NumOut](uint(v.Uint()))
	case reflect.Int8:
		//nolint:gosec // the int8 is confirmed
		return convertFromNumber[NumOut](int8(v.Int()))
	case reflect.Uint8:
		//nolint:gosec // the uint8 is confirmed
		return convertFromNumber[NumOut](uint8(v.Uint()))
	case reflect.Int16:
		//nolint:gosec // the int16 is confirmed
		return convertFromNumber[NumOut](int16(v.Int()))
	case reflect.Uint16:
		//nolint:gosec // the uint16 is confirmed
		return convertFromNumber[NumOut](uint16(v.Uint()))
	case reflect.Int32:
		//nolint:gosec // the int32 is confirmed
		return convertFromNumber[NumOut](int32(v.Int()))
	case reflect.Uint32:
		//nolint:gosec // the uint32 is confirmed
		return convertFromNumber[NumOut](uint32(v.Uint()))
	case reflect.Int64:
		return convertFromNumber[NumOut](int64(v.Int()))
	case reflect.Uint64:
		return convertFromNumber[NumOut](uint64(v.Uint()))
	case reflect.Float32:
		return convertFromNumber[NumOut](float32(v.Float()))
	case reflect.Float64:
		return convertFromNumber[NumOut](float64(v.Float()))
	case reflect.Bool:
		o := 0
		if v.Bool() {
			o = 1
		}
		return NumOut(o), nil
	case reflect.String:
		return convertFromString[NumOut](v.String())
	}

	return 0, errorHelper{
		err: fmt.Errorf("%w from %T", ErrUnsupportedConversion, orig),
	}
}

// MustConvert calls [Convert] to convert the value to the desired type, and panics if the conversion fails.
func MustConvert[NumOut Number, NumIn Input](orig NumIn) NumOut {
	converted, err := Convert[NumOut](orig)
	if err != nil {
		panic(err)
	}
	return converted
}

func convertFromNumber[NumOut Number, NumIn Number](orig NumIn) (converted NumOut, err error) {
	converted = NumOut(orig)

	// floats could be compared directly
	switch any(converted).(type) {
	case float64:
		// float64 cannot overflow, so we don't have to worry about it
		return converted, nil
	case float32:
		origFloat64, isFloat64 := any(orig).(float64)
		if !isFloat64 {
			// only float64 can overflow float32
			// everything else can be safely converted
			return converted, nil
		}

		// check boundary
		if math.Abs(origFloat64) < math.MaxFloat32 {
			// the value is within float32 range, there is no overflow
			return converted, nil
		}

		// TODO: check for numbers close to math.MaxFloat32

		boundary := getUpperBoundary(converted)
		errBoundary := ErrExceedMaximumValue
		if negative(orig) {
			boundary = getLowerBoundary(converted)
			errBoundary = ErrExceedMinimumValue
		}

		return 0, errorHelper{
			value:    orig,
			err:      errBoundary,
			boundary: boundary,
		}
	}

	errBoundary := ErrExceedMaximumValue
	boundary := getUpperBoundary(converted)
	if negative(orig) {
		errBoundary = ErrExceedMinimumValue
		boundary = getLowerBoundary(converted)
	}

	if !sameSign(orig, converted) {
		return 0, errorHelper{
			value:    orig,
			err:      errBoundary,
			boundary: boundary,
		}
	}

	// convert back to the original type
	cast := NumIn(converted)
	// and compare
	base := orig
	switch f := any(orig).(type) {
	case float64:
		base = NumIn(math.Trunc(f))
	case float32:
		base = NumIn(math.Trunc(float64(f)))
	}

	// exact match
	if cast == base {
		return converted, nil
	}

	return 0, errorHelper{
		value:    orig,
		err:      errBoundary,
		boundary: boundary,
	}
}

func convertFromString[NumOut Number](s string) (converted NumOut, err error) {
	s = strings.TrimSpace(s)

	if b, err := strconv.ParseBool(s); err == nil {
		if b {
			return NumOut(1), nil
		}
		return NumOut(0), nil
	}

	if strings.Contains(s, ".") {
		o, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, errorHelper{
				value: s,
				err:   fmt.Errorf("%w %v to %T", ErrStringConversion, s, converted),
			}
		}
		return convertFromNumber[NumOut](o)
	}

	if strings.HasPrefix(s, "-") {
		o, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return 0, errorHelper{
					value:    s,
					err:      ErrExceedMinimumValue,
					boundary: math.MinInt,
				}
			}
			return 0, errorHelper{
				value: s,
				err:   fmt.Errorf("%w %v to %T", ErrStringConversion, s, converted),
			}
		}

		return convertFromNumber[NumOut](o)
	}

	o, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			return 0, errorHelper{
				value:    s,
				err:      ErrExceedMaximumValue,
				boundary: uint(math.MaxUint),
			}
		}

		return 0, errorHelper{
			value: s,
			err:   fmt.Errorf("%w %v to %T", ErrStringConversion, s, converted),
		}
	}
	return convertFromNumber[NumOut](o)
}

// ToInt attempts to convert any [Type] value to an int.
// If the conversion results in a value outside the range of an int,
// an [ErrConversionIssue] error is returned.
func ToInt[T Number](i T) (int, error) {
	return convertFromNumber[int](i)
}

// ToUint attempts to convert any [Number] value to an uint.
// If the conversion results in a value outside the range of an uint,
// an [ErrConversionIssue] error is returned.
func ToUint[T Number](i T) (uint, error) {
	return convertFromNumber[uint](i)
}

// ToInt8 attempts to convert any [Number] value to an int8.
// If the conversion results in a value outside the range of an int8,
// an [ErrConversionIssue] error is returned.
func ToInt8[T Number](i T) (int8, error) {
	return convertFromNumber[int8](i)
}

// ToUint8 attempts to convert any [Number] value to an uint8.
// If the conversion results in a value outside the range of an uint8,
// an [ErrConversionIssue] error is returned.
func ToUint8[T Number](i T) (uint8, error) {
	return convertFromNumber[uint8](i)
}

// ToInt16 attempts to convert any [Number] value to an int16.
// If the conversion results in a value outside the range of an int16,
// an [ErrConversionIssue] error is returned.
func ToInt16[T Number](i T) (int16, error) {
	return convertFromNumber[int16](i)
}

// ToUint16 attempts to convert any [Number] value to an uint16.
// If the conversion results in a value outside the range of an uint16,
// an [ErrConversionIssue] error is returned.
func ToUint16[T Number](i T) (uint16, error) {
	return convertFromNumber[uint16](i)
}

// ToInt32 attempts to convert any [Number] value to an int32.
// If the conversion results in a value outside the range of an int32,
// an [ErrConversionIssue] error is returned.
func ToInt32[T Number](i T) (int32, error) {
	return convertFromNumber[int32](i)
}

// ToUint32 attempts to convert any [Number] value to an uint32.
// If the conversion results in a value outside the range of an uint32,
// an [ErrConversionIssue] error is returned.
func ToUint32[T Number](i T) (uint32, error) {
	return convertFromNumber[uint32](i)
}

// ToInt64 attempts to convert any [Number] value to an int64.
// If the conversion results in a value outside the range of an int64,
// an [ErrConversionIssue] error is returned.
func ToInt64[T Number](i T) (int64, error) {
	return convertFromNumber[int64](i)
}

// ToUint64 attempts to convert any [Number] value to an uint64.
// If the conversion results in a value outside the range of an uint64,
// an [ErrConversionIssue] error is returned.
func ToUint64[T Number](i T) (uint64, error) {
	return convertFromNumber[uint64](i)
}

// ToFloat32 attempts to convert any [Number] value to a float32.
// If the conversion results in a value outside the range of a float32,
// an [ErrConversionIssue] error is returned.
func ToFloat32[T Number](i T) (float32, error) {
	return convertFromNumber[float32](i)
}

// ToFloat64 attempts to convert any [Number] value to a float64.
// If the conversion results in a value outside the range of a float64,
// an [ErrConversionIssue] error is returned.
func ToFloat64[T Number](i T) (float64, error) {
	return convertFromNumber[float64](i)
}
