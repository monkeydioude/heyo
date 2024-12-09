package assert

import (
	"errors"
	"log"
)

var ErrErrorShouldBeNil = errors.New("'error' SHOULD be 'nil'")
var ErrPointerShouldNotBeNil = errors.New("'pointer' SHOULD NOT be 'nil'")
var ErrShouldNotBeEmpty = errors.New("value should not be empty")

func NoError(err error, wrap ...error) {
	if err != nil {
		log.Fatalf("%s: %s: %s", ErrErrorShouldBeNil, err, errors.Join(wrap...))
	}
}

func NotNil(ptr any, wrap ...error) {
	if ptr == nil {
		log.Fatalf("%s: %s", ErrErrorShouldBeNil, errors.Join(wrap...))
	}
}

func NotEmpty(value any, wrap ...error) {
	switch v := value.(type) {
	case nil:
		log.Fatalf("nil value: %s: %s: %s", ErrShouldNotBeEmpty, ErrPointerShouldNotBeNil, errors.Join(wrap...))
	case string:
		if v == "" {
			log.Fatalf("string value: %s: %s", ErrShouldNotBeEmpty, errors.Join(wrap...))
		}
	}
}

type AssertTrial interface {
	*string
	any
}

func NotNilNorEmpty[T AssertTrial](value T, wrap ...error) {
	NotNil(value, wrap...)
	NotEmpty(*value, wrap...)
}
