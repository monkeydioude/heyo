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

func NotNil(ptr any) {
	if ptr == nil {
		log.Fatal(ErrErrorShouldBeNil)
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
