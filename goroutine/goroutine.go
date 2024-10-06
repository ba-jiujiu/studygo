package goroutine

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	jobChan    chan int64
	resultChan chan int
	wg         sync.WaitGroup
)

func CalcDigSum() {
	defer wg.Done()
	for {
		n, ok := <-jobChan
		if !ok {
			break
		}
		sum := 0
		for n > 0 {
			sum += int(n % 10)
			n /= 10
		}
		resultChan <- sum
	}
}

func GetRandom(n int) {
	defer close(jobChan)
	for i := 0; i < n; i++ {
		r := rand.Int63()
		jobChan <- r
	}
}

func init() {
	jobChan = make(chan int64, 12)
	resultChan = make(chan int, 6)
}

func Do() {
	rand.Seed(time.Now().UnixNano())

	go GetRandom(24)
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go CalcDigSum()
	}
	// 开启这个线程的作用是当所有CalcDigSum执行完毕后关闭resultChan通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for v := range resultChan {
		fmt.Println(v)
	}
}
