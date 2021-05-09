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

func NewServiceNotFailingTranslator() *Service {
    t := newRandomTranslator(
        100*time.Millisecond,
        500*time.Millisecond,
        0,
    )


    return NewServiceFrom(t)
}

func NewServiceFrom(t Translator) *Service{
    return &Service{
        translator: t,
        cache: newMapCacheRepository(),
        dedublicate: newStandardDedublicateService(),
    }
}

func TestDeduplicateWithSameParametersNotPossibleWhenTheTranslateReturnsErrorForOneButNoErrorForTheOther(t *testing.T) {
    ctx := context.Background()
    commonCache := newMapCacheRepository()

    stubTranslator := newRandomTranslator(
        100*time.Millisecond,
        500*time.Millisecond,
        0,
    )

    failingTr := newRandomTranslator(
        100*time.Millisecond,
        500*time.Millisecond,
        1,
    )
    commonDedublicate := newStandardDedublicateService()

    service2 := &Service{
        translator: stubTranslator,
        cache: commonCache,
        dedublicate: commonDedublicate,
    }



    failingService := &Service{
        translator: failingTr,
        cache: commonCache,
        dedublicate: commonDedublicate,
    }
    failingService.cache = service2.cache
    if service2.translator.(*randomTranslator).errorProb !=0 {
        t.Errorf("Expect error prob 0, but got %f", service2.translator.(randomTranslator).errorProb)

    }

   // log.SetOutput(ioutil.Discard)
    defer log.SetOutput(os.Stderr)


    requestModel := RequestModel{
        language.Polish,
        language.Japanese,
        "Dzien dobry",
    }
    ctx1 := context.WithValue(ctx, "whereami", "goroutine1")
    ctx2 := context.WithValue(ctx, "whereami", "goroutine2")


    var result1, result2 string
    var err1, err2 error

    var wg sync.WaitGroup
    var wgGlobal sync.WaitGroup
    wg.Add(1)
    go func() {
        log.Println("FailingService is the first service to run")
        wg.Done()
        wgGlobal.Add(1)
        result1, err1 = failingService.TranslateOnce(ctx1, requestModel)
        wgGlobal.Done()
    }()
    wg.Wait()
    wgGlobal.Add(1)
    // failingService must be the first service2 to run
    go func() {
        log.Println("service2 is sequentially the second service to run")
        result2, err2 = service2.TranslateOnce(ctx2, requestModel)
        wgGlobal.Done()
    }()

    wgGlobal.Wait()
    if err1 == nil {
        t.Errorf("test setup error: expect First service to fail and result1 to be \"\", but got result1: %v", result1)
    }

    if result2 == "" || err2 != nil {
        t.Errorf("Expect result2 not to be empty, service2 must succeed after waiting for failing service to run its own `Translate`" +
            "result2: \"%s\" err2: \"%v\"", result2, err2)

    }
}


func TestDeduplicateServiceWhenTranslationServiceIsStable(t *testing.T) {
    ctx := context.Background()
    service := NewServiceNotFailingTranslator()
    log.SetOutput(ioutil.Discard)
    defer log.SetOutput(os.Stderr)

    ctx1 := context.WithValue(ctx, "whereami", "goroutine1")
    ctx2 := context.WithValue(ctx, "whereami", "goroutine2")

    howTranslate := RequestModel{
        language.English,
        language.Japanese,
        "Do you promise?",
    }


    t.Run("blackbox call simulatiously with same parameters", func(t *testing.T) {
        if service.dedublicate.BusyWith(howTranslate) {
            t.Errorf("Expect BusyWith(%v) be false by default", howTranslate)
        }

        var wg sync.WaitGroup
        wg.Add(1)
        var result1, result2 string
        go func() {
            result1, _ = service.TranslateOnce(ctx1, howTranslate)
            wg.Done()
        }()

        wg.Add(1)
        go func() {
            result2, _ = service.TranslateOnce(ctx2, howTranslate)
            wg.Done()
        }()
        time.Sleep(5 * time.Millisecond)
        if !service.dedublicate.BusyWith(howTranslate) {
            t.Errorf("Expect inProgress to be true while executing Translation(%v)"+
                "try increasing time.Sleep delay before checking", howTranslate)
        }
        wg.Wait()

        if service.dedublicate.BusyWith(howTranslate) {
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
        result1, err1 = s.TranslateOnce(ctx, useCase1)
        fmt.Println(result1, err1)
        wg.Done()
    }()

    wg.Add(1)
    go func() {
        result2, err2 = s.TranslateOnce(ctx, useCase2)
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

func TestDedublicateWhenRequestModelIsDifferentObjectButSameParametersInsideDeeply(t *testing.T) {}
func TestTwoServicesFailes(t *testing.T) {}
