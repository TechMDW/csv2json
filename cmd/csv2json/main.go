package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	csv2json "github.com/TechMDW/csv2json/pkg"
)

func main() {
	filepath := flag.String("file", "", "Path to the CSV file")
	seperator := flag.String("seperator", "", "CSV seperator, default is auto detect")
	header := flag.Bool("header", false, "CSV has header")

	flag.Parse()

	var sep rune
	if *seperator != "" {
		sep = []rune(*seperator)[0]
	}

	if *filepath != "" {
		data, err := csv2json.ParseFile(*filepath, sep)
		if err != nil {
			panic(err)
		}

		jsonData, err := data.ToJSON(*header)
		if err != nil {
			panic(err)
		}

		println(string(jsonData))
		return
	}

	// Workaround to avoid blocking (io.ReadAll) on stdin when no data input.
	// Read lines from stdin and append to csvData until EOF.
	lineChan := make(chan string)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lineChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			errChan <- err
		}
		close(doneChan)
	}()

	timeout := 1000 * time.Millisecond
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	var csvData []byte

	for {
		select {
		case line := <-lineChan:
			csvData = append(csvData, []byte(line+"\n")...)
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout)
		case err := <-errChan:
			panic(err)
		case <-doneChan:
			if len(csvData) == 0 {
				fmt.Println("No CSV data provided")
				return
			}
			data, err := csv2json.Parse(csvData, sep)
			if err != nil {
				panic(err)
			}

			jsonData, err := data.ToJSON(*header)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(jsonData))
			return
		case <-timer.C:
			fmt.Println("No CSV data provided (timeout)")
			return
		}
	}
}
