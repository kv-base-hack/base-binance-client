package server

import "errors"

var (
	ErrInvalidKLineRequest = errors.New("invalid kline request")
	ErrInvalidKLineLength  = errors.New("invalid kline length")

	ErrInvalidOIStatsRequest = errors.New("invalid OI stats")
	ErrInvalidBookTicker     = errors.New("invalid book ticker")
)
