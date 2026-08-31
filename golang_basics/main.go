package main

import (
	"context"
	"fmt"
	// "sync"
	// "sync/atomic"
)

func worker(ctx context.Context) {

}

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
// 	defer cancel()

// 	go worker(ctx)
// 	time.Sleep(400 * time.Millisecond)
// }

func main() {
	for i := range 5 {
		go func() {
			fmt.Println("hi", i)
		}()
	}
}

// atomic package
// func main() {
// 	var counter int64
// 	var wg sync.WaitGroup

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			atomic.AddInt64(&counter, 1)
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println("the counter value is", atomic.LoadInt64(&counter))
// }

// var config map[string]string
// var configOnce sync.Once

// func loadConfig() {
// 	fmt.Println("this is a expensive operation should happen once")
// 	config = map[string]string{"env": "production"}
// }

// func getConfig() map[string]string {
// 	configOnce.Do(loadConfig)
// 	return config
// }

// func main() {
// 	var wg sync.WaitGroup

// 	for range 10 {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			getConfig()
// 		}()
// 	}

// 	wg.Wait()
// }

// readwrite mutex
// type Cache struct {
// 	mu   sync.RWMutex
// 	data map[string]string
// }

// func newCache() *Cache {
// 	return &Cache{data: make(map[string]string)}
// }

// func (c *Cache) get(key string) (string, bool) {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()
// 	val, found := c.data[key]
// 	return val, found
// }

// func (c *Cache) set(key, val string) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.data[key] = val
// }

// func main() {
// 	cache := newCache()
// 	cache.set("name", "anuj")

// 	var wg sync.WaitGroup
// 	for range 10 {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			val, _ := cache.get("name")
// 			fmt.Println("the val is : ", val)
// 		}()
// 	}

// 	wg.Wait()
// }

// mutex code
// type BankAccount struct {
// 	mu      sync.Mutex
// 	balance int
// }

// func (b *BankAccount) deposit(amount int) {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	b.balance += amount
// }

// func (b *BankAccount) checkBalance() int {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	return b.balance
// }

// func main() {
// 	account := &BankAccount{}
// 	var wg sync.WaitGroup

// 	for range 100 {
// 		wg.Go(func() {
// 			account.deposit(1)
// 		})
// 	}

// 	wg.Wait()
// 	fmt.Println("the remaining balance is ", account.checkBalance())
// }

// waitgroup code
// func main() {
// 	var wg sync.WaitGroup

// 	tasks := []string{"hey ", "how ", "are ", "you ?"}
// 	for _, task := range tasks {
// 		wg.Add(1)
// 		go func(t string) {
// 			defer wg.Done()
// 			fmt.Println(t)
// 		}(task)
// 	}
// 	wg.Wait()
// 	fmt.Println("all task done")
// }

// func produce(ch chan<- int) {
// 	for i := 0; i < 5; i++ {
// 		ch <- i
// 	}
// 	close(ch)
// }

// func recieve(ch <-chan int) {
// 	for v := range ch {
// 		fmt.Println("value recieved from channel is : ", v)
// 	}
// }

// func main() {
// 	fastch := make(chan string)
// 	slowch := make(chan string)

// 	go func() {
// 		time.Sleep(100 * time.Millisecond)
// 		slowch <- "this is the slow channel"
// 	}()

// 	go func() {
// 		time.Sleep(11 * time.Millisecond)
// 		fastch <- "this is the fast ch"
// 	}()

// 	select {
// 	case msg := <-slowch:
// 		fmt.Println(msg)
// 	case msg := <-fastch:
// 		fmt.Println(msg)
// 	case <-time.After(10 * time.Millisecond):
// 		fmt.Println("timed out waiting : ")
// 	default:
// 		fmt.Println("taking too much time moving on")
// 	}
// }

// func main() {
// 	counter := 0
// 	var wg sync.WaitGroup

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			counter++
// 			fmt.Println("The value of id is : ", i)
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println("final counter: ", counter)
// }

// func leaky() <-chan int {
// 	ch := make(chan int)
// 	go func() {
// 		ch <- 42
// 	}()
// 	return ch
// }

// func main() {
// 	leaky()
// 	fmt.Println("hi main is finsihed")
// }
