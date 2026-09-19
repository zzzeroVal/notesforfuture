package main

import (
	"fmt"
	"strconv"
	"strings"
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
	// inputData := []int{0, 1}
	data := make([]int, 0)

	for i, val := range data {
		newVal := strconv.Itoa(val)

		md5 := DataSignerMd5(newVal)
		crc32Data := DataSignerCrc32(newVal)
		crc32md5 := DataSignerCrc32(DataSignerMd5(newVal))
		result := DataSignerCrc32(newVal) + "~" + DataSignerCrc32(DataSignerMd5(newVal))

		fmt.Printf("%d SingleHash data %s\n", i, newVal)
		fmt.Printf("%d SingleHash md5(data) %s", i, md5)
		fmt.Printf("%d SingleHash crc32(md5(data)) %s", i, crc32md5)
		fmt.Printf("%d SingleHash crc32(data) %s", i, crc32Data)
		fmt.Printf("%d SingleHash result %s", i, result)
	}
}

func MultiHash(in, out chan interface{}) {
	data := make([]int, 0)
	th := []int{0, 1, 2, 3, 4, 5}
	var multiHashResult strings.Builder
	for _, val := range data {
		newVal := strconv.Itoa(val)
		result := DataSignerCrc32(newVal) + "~" + DataSignerCrc32(DataSignerMd5(newVal))
		for _, thNum := range th {
			thNumStr := strconv.Itoa(thNum)
			crc32 := DataSignerCrc32(thNumStr + newVal)
			fmt.Printf("%s MultiHash: crc32(th+step1)) %s %s\n", result, thNumStr, crc32)
			multiHashResult.WriteString(crc32)
		}
		fmt.Printf("%s MultiHash result: %s", result, multiHashResult.String())
	}
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
