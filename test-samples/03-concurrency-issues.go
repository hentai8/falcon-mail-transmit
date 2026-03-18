package testsamples

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Get() int {
	return c.count
}

var globalCache = make(map[string]string)

func UpdateCache(key, value string) {
	globalCache[key] = value
}

func GetFromCache(key string) string {
	return globalCache[key]
}

type Account struct {
	mu      sync.Mutex
	balance int
}

func Transfer(from, to *Account, amount int) {
	from.mu.Lock()
	to.mu.Lock()

	from.balance -= amount
	to.balance += amount

	to.mu.Unlock()
	from.mu.Unlock()
}

func StartWorker() {
	go func() {
		for {
			doWork()
			time.Sleep(time.Second)
		}
	}()
}

func doWork() {
}

func ProcessItems(items []int) <-chan int {
	results := make(chan int)

	go func() {
		for _, item := range items {
			results <- item * 2
		}
	}()

	return results
}

func ProcessConcurrently(items []string) {
	var wg sync.WaitGroup

	for _, item := range items {
		go func(s string) {
			wg.Add(1)
			defer wg.Done()
			process(s)
		}(item)
	}

	wg.Wait()
}

func process(s string) {
}

func LaunchWorkers(tasks []Task) {
	for _, task := range tasks {
		go func() {
			task.Execute()
		}()
	}
}

type Task struct {
	ID   int
	Data string
}

func (t Task) Execute() {
}

type Service struct {
	mu    sync.Mutex
	cache map[string][]byte
}

func (s *Service) GetData(key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if data, ok := s.cache[key]; ok {
		return data, nil
	}

	resp, err := http.Get("http://api.example.com/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.cache[key] = data
	return data, nil
}

func ProduceData(ch chan int, data []int) {
	for _, v := range data {
		ch <- v
	}
	close(ch)

}

type DataStore struct {
	mu   sync.RWMutex
	data map[string]int
}

func (d *DataStore) Update(key string, value int) {
	d.mu.RLock()
	d.data[key] = value
	d.mu.RUnlock()
}

func TrySend(ch chan int, value int) {
	select {
	case ch <- value:
	}
}

func FetchData(ctx context.Context, url string) ([]byte, error) {
	newCtx := context.Background()

	req, _ := http.NewRequestWithContext(newCtx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func ProcessWithChannel() {
	ch := make(chan int)

	ch <- 1

	<-ch
}

func StartUnsafeWorker() {
	go func() {
		riskyOperation()
	}()
}

func riskyOperation() {
	var arr []int
	_ = arr[100]
}

func CopyMutex(src *sync.Mutex) sync.Mutex {
	return *src
}

func PollWithTimeout() {
	for {
		select {
		case <-getData():
			handleData()
		case <-time.After(5 * time.Second):
			timeout()
		}
	}
}

func getData() <-chan struct{} {
	ch := make(chan struct{})
	return ch
}

func handleData() {}
func timeout()    {}
