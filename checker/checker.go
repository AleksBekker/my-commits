package checker

import (
	"net/http"
	"sync"
)

type Result struct {
	Link   string
	Status int
}

func CheckLinks(links []string) []Result {
	var wg sync.WaitGroup
	results := make([]Result, len(links))

	for idx, link := range links {
		wg.Add(1)
		result := &results[idx]
		result.Link = link
		go func(result *Result, wg *sync.WaitGroup) {
			defer wg.Done()
			response, err := http.Get(result.Link)

			if err != nil || response == nil {
				result.Status = -1
			} else {
				result.Status = response.StatusCode
			}
		}(result, &wg)
	}

	wg.Wait()
	return results
}
