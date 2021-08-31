package main


import (
        "sync"
	"strconv"
	"sort"
	"strings"
	"time"
)


func ExecutePipeline(jobs ...job) {
	ch := make([]chan interface{}, len(jobs)+1)
	for i := range ch {
		ch[i] = make(chan interface{}, 100)
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(jobs))
	for idx, j := range jobs {
		go func(in chan interface{}, out chan interface{}, j job) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(ch[idx], ch[idx+1], j)
	}
	wg.Wait()
}


func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	wg.Add(7)
	for i := 0; i < 7; i++ {
		time.Sleep(10 * time.Millisecond)
		data := (<-in).(int)
		go func(data int, wg *sync.WaitGroup) {
			defer wg.Done()
			w := &sync.WaitGroup{}
			w.Add(1)
			first := make(chan string, 1)
			go func(data int, first chan<- string, w *sync.WaitGroup) {
				defer w.Done()
				first <- DataSignerCrc32(strconv.Itoa(data))
			}(data, first, w)
			second := DataSignerCrc32(DataSignerMd5(strconv.Itoa(data)))
		        w.Wait()
			f := <-first
		        out <- f + "~" + second
		}(data, wg)
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}


func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	wg.Add(7)
	for i := 0; i < 7; i++ {
		go func(in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			data := (<-in).(string)
			res := make([]string, 6)
			w := &sync.WaitGroup{}
			w.Add(6)
			for i := 0; i < 6; i++ {
				go func(res []string, idx int, w *sync.WaitGroup) {
					defer w.Done()
					res[idx] = DataSignerCrc32(strconv.Itoa(idx) + data)
				}(res, i, w)
			}
			w.Wait()
			out <- strings.Join(res, "")
		}(in, out, wg)
	}
	wg.Wait()
}


func CombineResults(in, out chan interface{}) {
	sl := make([]string, 0)
	for input := range in {
		str := input.(string)
		sl = append(sl, str)
	}
	sort.Strings(sl)
	out <- strings.Join(sl, "_")
}
