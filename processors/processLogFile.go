package processors

import (
	"fmt"
	"io"
	"mime/multipart"
	"sort"
	"strings"
	"sync"
)

type Output struct {
	Key   string
	Value int
}

func ProcessLogFile(file multipart.File, keywords []string) ([]Output, error) {
	chunkSize := 4096
	defer file.Close()

	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("Error determining file size: %v\n", err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("Error resetting file pointer: %v\n", err)
	}

	results := make(chan map[string]int)
	var wg sync.WaitGroup

	for offset := int64(0); offset < size; offset += int64(chunkSize) {
		wg.Add(1)
		go CountKeywords(file, offset, chunkSize, keywords, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalCounts := make(map[string]int)
	for result := range results {
		if result == nil {
			continue
		}
		for keyword, count := range result {
			totalCounts[keyword] += count
		}
	}

	sortedMap := sortMapDescending(totalCounts)
	return sortedMap, nil
}

func CountKeywords(file multipart.File, offset int64, chunkSize int, keywords []string, results chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	buffer := make([]byte, chunkSize)

	_, err := file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading chunk: %v\n", err)
		results <- nil
		return
	}

	countMap := make(map[string]int)
	chunkStr := string(buffer)

	for _, keyword := range keywords {
		upperCaseKeyWord := strings.ToUpper(keyword)
		countMap[upperCaseKeyWord] = strings.Count(strings.ToUpper(chunkStr), upperCaseKeyWord)
	}

	results <- countMap
}

func sortMapDescending(m map[string]int) []Output {
	sortedMap := make([]Output, 0, len(m))

	for k, v := range m {
		sortedMap = append(sortedMap, Output{Key: k, Value: v})
	}

	sort.Slice(sortedMap, func(i, j int) bool {
		return sortedMap[i].Value > sortedMap[j].Value
	})
	return sortedMap
}
