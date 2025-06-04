package repeatible

import (
	"time"
)

func DoWithTries(fn func() error, attemps int, delay time.Duration) (err error) {
	for attemps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemps--

			continue
		}

		return nil
	}

	return
}
