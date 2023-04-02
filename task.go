package taskflow

import (
	"bufio"
	"log"
	"os"
	"sync"
)

type Factory interface {
	Make(line string) Task
}

type Task interface {
	Process()
	Print()
}

func Run(f Factory, workers ...int) {
	defaultWorkers := 1000
	if len(workers) > 0 {
		defaultWorkers = workers[0]
	}

	var wg sync.WaitGroup

	in := make(chan Task)

	wg.Add(1)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			in <- f.Make(s.Text())
		}
		if s.Err() != nil {
			log.Fatalf("error reading stdin: %s", s.Err())
		}
		close(in)
		wg.Done()
	}()

	out := make(chan Task)

	for i := 0; i < defaultWorkers; i++ {
		wg.Add(1)
		go func() {
			for t := range in {
				t.Process()
				out <- t
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for t := range out {
		t.Print()
	}
}
