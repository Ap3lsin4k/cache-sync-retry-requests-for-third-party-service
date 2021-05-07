package main_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/pailcamper/pc-offline-challenge"
	"golang.org/x/text/language"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type DummyService struct {
	translator translatorSpy

}

func NewDService() *DummyService {
	t := newTranslatorSpy(
		100*time.Millisecond,
		500*time.Millisecond,
		1,
	)

	return &DummyService{
		translator: *t,
	}
}

// translatorSpy in a Translator implementation which is used
// only for testing purposes
type translatorSpy struct {
	minDelay      time.Duration
	maxDelay      time.Duration
	errorProb     float64
	CalledCounter uint8
}

func newTranslatorSpy(minDelay, maxDelay time.Duration, errorProbability float64) *translatorSpy {
	return &translatorSpy{
		minDelay:      minDelay,
		maxDelay:      maxDelay,
		errorProb:     errorProbability,
		CalledCounter: 0,
	}
}

func (t *translatorSpy) Translate(_ context.Context, from, to language.Tag, data string) (string, error) {
//	time.Sleep(t.randomDuration())
	t.CalledCounter++

	if rand.Float64() < t.errorProb {
		return "", errors.New("translation failed")
	}

	res := fmt.Sprintf("%v -> %v : %v -> %v", from, to, data, strconv.FormatInt(rand.Int63(), 10))
	return res, nil
}

func (t translatorSpy) randomDuration() time.Duration {
	delta := t.maxDelay - t.minDelay
	var delay time.Duration = t.minDelay + time.Duration(rand.Int63n(int64(delta)))
	return delay
}

func (t *translatorSpy) IncrementCalledCounter() {
	t.CalledCounter += 1
}


func f(s *DummyService) func() (string, error) {
	ctx := context.Background()
	return func () (string, error) {
		return s.translator.Translate(ctx, language.English, language.Japanese, "test")
	}
}

func TestNextDelayBeforeConnectingIsHigher(t *testing.T) {
	t.Run("how much wait before connecting again", func(t *testing.T) {
		previous, _ := main.NextBackoff(0)
		next, _ := main.NextBackoff(1)

		if previous >= next {
			t.Errorf("Next backoff must be greater than previous")
		}

		got, _ := main.NextBackoff(9)
		if got > 15 * time.Minute {
			t.Errorf("Last backoff before panicking should not be greater than 10 minutes, but got %d minutes", got/time.Minute)
		}

		_, err := main.NextBackoff(10)
		if err == nil{
			t.Errorf("Should give up when number of attempts to connect to the server is greater than threshhold")
		}
	})

	s := NewDService()

	translate :=  f(s)
	fmt.Println("<")

	delayBeforeReconnecting := func (uint8) (time.Duration, error) {
		return 100*time.Millisecond, nil
	}

	expectAttemptsBeforePanicking := uint8(4)
	_, _ = main.TryTranslate(translate, delayBeforeReconnecting, expectAttemptsBeforePanicking)

	mySpy := s.translator
	got := s.translator.CalledCounter
	if got != expectAttemptsBeforePanicking {
		t.Errorf("Expect %d attempts before panicking, but got %d", 4, mySpy.CalledCounter)
	}
}

func g(translator *translatorSpy) func() string {
	return func () string {
		translator.IncrementCalledCounter()
		return "sc"
	}
}

func TestNothing(t *testing.T) {
	tr := newTranslatorSpy(
		100*time.Millisecond,
		500*time.Millisecond,
		1,
	)
	translate := g(tr)
	translate()
	if tr.CalledCounter < 1 {
		t.Errorf("Expect CalledCounter increased, but got %d", tr.CalledCounter)
	}
	translate()
	if tr.CalledCounter != 2 {
		t.Errorf("Expect CalledCounter increased, but got %d", tr.CalledCounter)
	}
}