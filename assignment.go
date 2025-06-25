package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	inputFiles := []string{"server1.log", "server2.log", "server3.log"}
	err := LogProcess(inputFiles, "error.log")
	if err != nil {
		log.Println("Error occured", err)
	} else {
		log.Println("Successfully logged in error log")
	}
}

func LogProcess(inputFiles []string, out string) error {
	errChan := make(chan string, 100)
	wg := &sync.WaitGroup{}

	writerWg := &sync.WaitGroup{}
	writerWg.Add(1)

	go func() {
		defer writerWg.Done()
		writeFiles(out, errChan)
	}()

	for _, file := range inputFiles {
		wg.Add(1)
		go readFiles(file, errChan, wg)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	writerWg.Wait()

	return nil
}

func readFiles(files string, errChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(files)
	if err != nil {
		log.Println("Error occured in opening a file", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ERROR") {
			errChan <- line
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Scanner error occured", err)
	}
}

func writeFiles(out string, errChan chan string) {
	file, err := os.Create(out)
	if err != nil {
		log.Println("Error occured in creating file", err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	for line := range errChan {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Println("Error writing on file", err)
		}
	}

	writer.Flush()
}
