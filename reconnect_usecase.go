package main

import (
	"errors"
	"fmt"
	"math"
	"time"
)

func TryTranslateOrRetry(translate func() (string, error), delayBeforeReconnecting func(uint8) (time.Duration, error), attemptsLeftBeforePanicking uint8) (string, error) {
	err := fmt.Errorf("attempted %d times but failed to connect", attemptsLeftBeforePanicking)
	for attempts := uint8(0); attempts <= attemptsLeftBeforePanicking+2; attempts++ {
		var translated = ""
		translated, err = translate()

		if err != nil {

			backOff, backOffErr := delayBeforeReconnecting(attempts)
			if backOffErr != nil {
				return "", backOffErr
			}

			// todo logger: fmt.Printf("%s. Connecting again in %d seconds...\n", err, backOff/time.Second)
			time.Sleep(backOff)
			attempts++
		} else {
			return translated, nil
		}
	}
	return "", err
}

func NextBackoff(connectingAttempts int) (time.Duration, error) {
	threshold := 10 * time.Minute
	backoff := nextBackoff(connectingAttempts)
	if backoff < threshold {
		return backoff, nil
	}
	return backoff, errors.New("backoff time before reconnecting is too long")
}

func nextBackoff(connectingAttempts int) time.Duration {
	return 100 * time.Millisecond * time.Duration(math.Pow(math.E, float64(connectingAttempts)))
}

