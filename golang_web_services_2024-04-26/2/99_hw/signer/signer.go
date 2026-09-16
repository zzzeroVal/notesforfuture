package main

import (
	"fmt"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	wg := &sync.WaitGroup{}

	for i, smthJob := range jobs {
		wg.Add(1)
		fmt.Println("Создали out", i+1)
		out := make(chan interface{})

		go func(in, out chan interface{}, wg *sync.WaitGroup, smthJob job, i int) {
			defer wg.Done()
			fmt.Println("Запустили job", i+1)
			smthJob(in, out)
			close(out)
			fmt.Println("Закрыли out", i+1)
		}(in, out, wg, smthJob, i)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
}

func MultiHash(in, out chan interface{}) {
}

func CombineResults(in, out chan interface{}) {
}

/*func main() {

	job1 := func(in, out chan interface{}) {
		value := <-in
		fmt.Println("job1 получила:", value)
		result := value.(int) * 2
		fmt.Println("job1 отправила:", result)
		out <- result
	}
	job2 := func(in, out chan interface{}) {
		value := <-in
		fmt.Println("job2 получила:", value)
		result := value.(int) + 10
		fmt.Println("job2 отправила:", result)
		out <- result
	}
	ExecutePipeline(job1, job2)
time.Sleep(time.Second)
}*/
