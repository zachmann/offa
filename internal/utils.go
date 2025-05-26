package internal

import (
	"crypto/rand"
	"encoding/base64"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// SliceIsSubsetOf checks if a slice `subset` is a subset of the slice `of`,
// i.e. it is verified that `of` contains all the elements of `subset`
func SliceIsSubsetOf[T comparable](subset []T, of []T) bool {
	for _, sub := range subset {
		if !slices.Contains(of, sub) {
			return false
		}
	}
	return true
}

// FirstNonEmpty is a utility function returning the first of the passed values that is not empty/zero
func FirstNonEmpty[C comparable](possibleValues ...C) C {
	var nullValue C
	for _, v := range possibleValues {
		if v != nullValue {
			return v
		}
	}
	return nullValue
}

// FirstNonEmptyFnc is a utility function returning the first of the passed
// values that is not empty/zero.
// In this function the values are not passed directly but a function that
// returns the value is passed instead. This enables lazy evaluation
func FirstNonEmptyFnc[C comparable](possibleValues ...func() C) C {
	var nullValue C
	for _, fnc := range possibleValues {
		if v := fnc(); v != nullValue {
			return v
		}
	}
	return nullValue
}

// FirstNonEmptyQueryParameter check a fiber.
// Ctx for multiple parameters and returns the value for the first one that is
// set
func FirstNonEmptyQueryParameter(c *fiber.Ctx, parameters ...string) string {
	var fncs []func() string
	for _, param := range parameters {
		fncs = append(
			fncs, func() string {
				return c.Query(param)
			},
		)
	}
	return FirstNonEmptyFnc(fncs...)
}

func RandomString(n int) (string, error) {
	byteLen := n * 3 / 4 // base64 expands by 4/3
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n], nil
}
