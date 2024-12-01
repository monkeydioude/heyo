package client

import "errors"

var ErrSubscriptionFailed = errors.New("subscription failed")
var ErrMessageReception = errors.New("message reception failed")
var ErrListenFatalErr = errors.New("fatal error during Listen")
var ErrNilPointer = errors.New("nil pointer")

var ErrHandshake = errors.New("could not perform handshake with the server")
