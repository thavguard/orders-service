package db

import (
	"database/sql"
	"errors"
	"strings"
)

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone) {
		return true
	}

	msg := strings.ToLower(err.Error())
	retryableErrors := []string{"deadlock", "timeout", "connection refused", "temporarily unavailable", "too many connections", "failed to connect"}

	for _, retryableErr := range retryableErrors {
		if strings.Contains(msg, retryableErr) {
			return true
		}
	}
	return false
}
