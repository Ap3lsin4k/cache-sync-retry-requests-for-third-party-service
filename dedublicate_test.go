package main

import (
    "context"
    "fmt"
    "golang.org/x/text/language"
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
    }
}

func TestDeduplicate(t *testing.T) {
    ctx := context.Background()
    service := NewStubService()

    if service.busy {
        t.Errorf("Expect inProgress be false by default")
    }

    fmt.Println("Before `    defer nextPart()`")
    defer nextPart(t, *service)
    fmt.Println("After `defer nextPart()`")

    _, _ = service.translator.Translate(ctx, language.English, language.Japanese, "test") //.promise(nextpart)
    if !service.busy {
        t.Errorf("Expect inProgress when called Translation")
    }
}

func nextPart(t *testing.T, service Service){
    if service.busy {
        t.Errorf("Expect inProgress be false after Translate finished")
    }
}