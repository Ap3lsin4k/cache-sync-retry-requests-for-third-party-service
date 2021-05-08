package main

import (
    "context"
    "fmt"
    "golang.org/x/text/language"
    "io/ioutil"
    "log"
    "os"
    "sync"
    "testing"
    "time"
)

func NewStubService() *Service {
    t := newRandomTranslator(
        100*time.Millisecond,
        500*time.Millisecond,
        0.1,
    )


    return &Service{
        translator: t,
        cache: map[RequestModel]string{},
        translatingInProgress: make(setOfRequestModel),
    }
}

func NewServiceNotFailingTranslator() *Service {
    t := newRandomTranslator(
        100*time.Millisecond,
        500*time.Millisecond,
        0,
    )


    return &Service{
        translator: t,
        cache: map[RequestModel]string{},
        translatingInProgress: make(setOfRequestModel),
    }
}



func TestDeduplicateForSame(t *testing.T) {
}
func TestWhiteboxServiceKnowsThatTranslationServiceIsExecuting(t *testing.T) {
    ctx := context.Background()
    service := NewStubService()
    log.SetOutput(ioutil.Discard)
    defer log.SetOutput(os.Stderr)


    howTranslate := RequestModel{
        language.English,
        language.Japanese,
        "Do you promise?",
    }

    useCase2 := RequestModel{
        from:       language.Ukrainian,
        to:         language.English,
        fromPhrase: "Тарас Шевченко",
    }


    t.Run("blackbox call simulatiously with same parameters", func(t *testing.T) {
        if service.BusyWith(howTranslate) {
            t.Errorf("Expect BusyWith(%v) be false by default", howTranslate)
        }
        if service.BusyWith(useCase2) {
            t.Errorf("Expect BusyWith(%v) be false by default", useCase2)
        }

        var wg sync.WaitGroup
        wg.Add(1)
        var result1, result2 string
        go func() {
            ctx1 := context.WithValue(ctx, "whereami", "goroutine1")
            result1, _ = service.TranslatePromise(ctx1, howTranslate)
            wg.Done()
        }()

        wg.Add(1)
        go func() {
            ctx2 := context.WithValue(ctx, "whereami", "goroutine2")
            result2, _ = service.TranslatePromise(ctx2, howTranslate)
            wg.Done()
        }()
        time.Sleep(100 * time.Millisecond)
        if !service.BusyWith(howTranslate) {
            t.Errorf("Expect inProgress to be true while executing Translation(%v)"+
                "try increasing time.Sleep delay before checking", howTranslate)
        }
        wg.Wait()

        if service.BusyWith(howTranslate) {
            t.Errorf("Expect inProgress to be false on fished Translation")
        }

        if result1 == "" || result1 != result2{
            t.Errorf("Expect results to be the same, but got result1 \"%s\" != result2 \"%s\"",
                result1, result2)
        }

    })
}

func TestNotDedublicateSimultaneousQueriesForDifferentParameters(t *testing.T) {

    ctx := context.Background()
    s := NewServiceNotFailingTranslator()
    log.SetOutput(ioutil.Discard)
    defer log.SetOutput(os.Stderr)
//  TODO Implement Logger in Base class as cross cutting concern: Aspect-oriented programming

    useCase1 := RequestModel{
        from:       language.English,
        to:         language.Japanese,
        fromPhrase: "same words",
    }
    useCase2 := RequestModel{
        from:       language.English,
        to:         language.Chinese,
        fromPhrase: "same words",
    }

    var result1, result2 string
    var err1, err2 error
    //:= "", "", nil, nil

    var wg sync.WaitGroup
    wg.Add(1)
    go func () {
        result1, err1 = s.TranslatePromise(ctx, useCase1)
        fmt.Println(result1, err1)
        wg.Done()
    }()

    wg.Add(1)
    go func() {
        result2, err2 = s.TranslatePromise(ctx, useCase2)
        wg.Done()
    }()

    wg.Wait()

    if result1 == "" || result2 == "" || err1 != nil || err2 != nil {
        t.Errorf("result1: %s; result2: %s; err1: %v; err2: %v;" +
            "this test should not check for error handling " +
            "please set `errorProbability` in fake translator to 0 " +
            "or wait for functions to finish",
            result1, result2, err1, err2)
    }

    if result1 == result2 {
        t.Errorf("Expect output to differ, but got %s==%s", result1, result2)
    }
}