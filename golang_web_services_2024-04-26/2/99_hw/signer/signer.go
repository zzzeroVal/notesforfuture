package main

import (
	"fmt"
	"sort"
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

func SingleHash(in chan interface{}, out chan interface{}) {
	valuesWG := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for val := range in {
		valuesWG.Add(1)
		newVal := strconv.Itoa(val.(int))
		go func(newVal string, val interface{}) {
			defer valuesWG.Done()
			var crc32Data string
			crc32WG := &sync.WaitGroup{}
			crc32WG.Add(1)

			go func() {
				defer crc32WG.Done()
				crc32Data = DataSignerCrc32(newVal)
			}()
			mu.Lock()
			md5 := DataSignerMd5(newVal)
			mu.Unlock()
			crc32md5 := DataSignerCrc32(md5)
			crc32WG.Wait()
			result := crc32Data + "~" + crc32md5

			out <- result
			fmt.Printf("%d SingleHash data %s\n", val, newVal)
			mu.Lock()
			fmt.Printf("%d SingleHash md5(data) %s\n", val, md5)
			mu.Unlock()
			fmt.Printf("%d SingleHash crc32(md5(data)) %s\n", val, crc32md5)
			fmt.Printf("%d SingleHash crc32(data) %s\n", val, crc32Data)
			fmt.Printf("%d SingleHash result %s\n", val, result)
		}(newVal, val)
	}
	valuesWG.Wait()
}

func MultiHash(in, out chan interface{}) {
	valuesWG := &sync.WaitGroup{}
	for val := range in {
		valuesWG.Add(1)
		go func(newVal string) {
			defer valuesWG.Done()
			wg := &sync.WaitGroup{}
			var multiHashResult strings.Builder
			results := make([]string, 6)

			for thNum := 0; thNum < 6; thNum++ {
				wg.Add(1)
				go func(thNum int, newVal string) {
					defer wg.Done()
					thNumStr := strconv.Itoa(thNum)
					crc32 := DataSignerCrc32(thNumStr + newVal)
					fmt.Printf("%s MultiHash: crc32(th+step1)) %s %s\n", newVal, thNumStr, crc32)
					results[thNum] = crc32
				}(thNum, newVal)
			}
			wg.Wait()
			for thNum := 0; thNum < 6; thNum++ {
				multiHashResult.WriteString(results[thNum])
			}
			out <- multiHashResult.String()
			fmt.Printf("%s MultiHash result: %s\n", newVal, multiHashResult.String())
		}(val.(string))
	}
	valuesWG.Wait()
}

func CombineResults(in, out chan interface{}) {
	var interimResult []string
	for val := range in {
		newVal := val.(string)
		interimResult = append(interimResult, newVal)
	}
	sort.Strings(interimResult)
	finalResult := strings.Join(interimResult, "_")
	out <- finalResult
	fmt.Printf("CombineResults %s\n", finalResult)
}
