package main

import "errors"

var ErrInMessageWasNil = errors.New("message to send was nil")
var ErrEnqueuingMessage = errors.New("error trying to enqueue a message")
var ErrAckNotOk = errors.New("ack code not ok")
var ErrMessageWasMalformed = errors.New("message was malformed, should be '@<event> <message>'")
var ErrMessageMinimum2Words = errors.New("message requires at least 2 words")
var ErrLeadingAtMissing = errors.New("event requires a leading @")
var ErrUpstreamClosed = errors.New("upstream closed")
var ErrSubscriptionFailed = errors.New("subscription failed")
var ErrMessageReception = errors.New("message reception failed")
var ErrReadingInput = errors.New("could not read input")
