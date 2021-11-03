package thread

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAll1(t *testing.T) {
	p1 := NewPromise(func(resolve Resolve) {
		tt := time.NewTimer(time.Second)
		<-tt.C

		fmt.Println("resolve p1")
		resolve.R(1)
	})

	p2 := NewPromise(func(resolve Resolve) {
		tt := time.NewTimer(time.Second)
		<-tt.C

		fmt.Println("resolve p2")
		resolve.R(2)
	})

	p3 := NewPromise(func(resolve Resolve) {
		tt := time.NewTimer(time.Second)
		<-tt.C

		fmt.Println("resolve p3")
		resolve.R(3)
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	All(p1, p2, p3).Then(func(result Result) {
		defer wg.Done()
		if result.Err != nil {
			fmt.Println("err: ", result.Err)
			return
		}

		fmt.Printf("success: %+v\n", result.Val)
	})

	wg.Wait()
}

func TestAll2(t *testing.T) {
	p1 := NewPromise(func(resolve Resolve) {
		tt := time.NewTimer(time.Second)
		<-tt.C

		fmt.Println("reject p1")
		resolve.E(errors.New("reject"))
	})

	p2 := NewPromise(func(resolve Resolve) {
		tt := time.NewTimer(time.Second)
		<-tt.C

		fmt.Println("reject p2")
		resolve.E(errors.New("reject"))
	})

	p3 := NewPromise(func(resolve Resolve) {
		tt := time.NewTimer(time.Second)
		<-tt.C

		fmt.Println("resolve p3")
		resolve.R(3)
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	All(p3, p1, p2).Then(func(result Result) {
		defer wg.Done()
		if result.Err != nil {
			fmt.Println("err: ", result.Err)
			return
		}

		fmt.Printf("success: %+v\n", result.Val)
	})

	wg.Wait()
}
