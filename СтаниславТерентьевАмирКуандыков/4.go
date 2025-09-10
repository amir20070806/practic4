package main

import (
 "fmt"
 "net/http"
 "sync"
 "time"
)

func main() {
 urls := []string{
  "https://youtube.com",
  "https://facebook.com",
  "https://github.com",
  "https://amazon.com",
 }

 jobs := make(chan string, len(urls))
 results := make(chan string, len(urls))

 var ff sync.WaitGroup

 for i := 1; i <= 3; i++ {
  ff.Add(1)
  go func(workerID int) {
   defer ff.Done()

   for url := range jobs {
    fmt.Printf("Воркер %d проверяет: %s\n", workerID, url)
    

    client := http.Client{
     Timeout: 10 * time.Second,
    }

    resp, err := client.Get(url)
    var status string
    if err != nil {
     status = fmt.Sprintf("ОШИБКА: %v", err)
    } else {
     defer resp.Body.Close()
     status = fmt.Sprintf("Статус: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
    }
    

    results <- fmt.Sprintf("Воркер %d: %s - %s", workerID, url, status)
   }
  }(i) 
 }

 for _, url := range urls {
  jobs <- url
 }
 close(jobs)


 ff.Wait()
 close(results)

 fmt.Println("\n URL:")
 for result := range results {
  fmt.Println(result)
 }
}
