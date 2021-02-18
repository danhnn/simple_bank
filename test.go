package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	results := make(chan int)
	rand.Seed(time.Now().UnixNano())
	j := 0

	for i := 0; i < 5; i++ {
		go func() {
			j++
			fmt.Println(j)
			time.Sleep(time.Duration(j) * time.Second)
			results <- rand.Int()
		}()
	}

	fmt.Println("HERE")
	for i := 0; i < 5; i++ {
		fmt.Println("HERE 2" + time.Now().String())
		res := <-results
		fmt.Println(res)
	}
}
