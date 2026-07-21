package worker

import "sync"

func Run(totalJobs int, workers int, job func(int)) {

	jobs := make(chan int)

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for i := range jobs {

				job(i)

			}

		}()

	}

	for i := 0; i < totalJobs; i++ {

		jobs <- i

	}

	close(jobs)

	wg.Wait()

}
