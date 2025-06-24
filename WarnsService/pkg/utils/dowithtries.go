package repeatible

import (
	"fmt"
	"time"
)

func DoWithTries(fn func() error, attemps int, delay time.Duration) (err error) {
	for attemps > 0 {
		if err = fn(); err != nil {
			attemps--

			time.Sleep(delay)
			continue
		}

		return nil
	}

	return fmt.Errorf("error after %d attemps got fail", attemps)
}
