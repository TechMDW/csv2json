package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/TechMDW/randish"
)

func main() {
	const numHeaders = 1000
	const numRows = 10000

	file, err := os.Create("test.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := make([]string, numHeaders)
	for i := 0; i < numHeaders; i++ {
		headers[i] = "Column" + strconv.Itoa(i+1)
	}
	if err := writer.Write(headers); err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}

	rand := randish.RandS()

	for i := 0; i < numRows; i++ {
		row := make([]string, numHeaders)
		for j := 0; j < numHeaders; j++ {
			row[j] = strconv.Itoa(rand.Intn(100) + 1)
		}
		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing row:", err)
			return
		}
	}

	fmt.Println("CSV file generated successfully.")
}
