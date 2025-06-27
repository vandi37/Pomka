package utils

import (
	"fmt"
	"time"

	e "errorspomka"
)

func DoWithTries(fn func() error, attemps int, delay time.Duration) (err error) {
	for i := attemps; i > 0; i-- {
		if err = fn(); err != nil {
			time.Sleep(delay)
			continue
		}

		return nil
	}

	return fmt.Errorf(e.ErrDoWithTries.Error(), attemps)
}
