package csv2json

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type CSVData [][]interface{}

func (c CSVData) ToJSON(header bool) ([]byte, error) {
	return ConvertToJSON(c, header)
}

// parseScanner parses the CSV data from a bufio.Scanner
func parseScanner(scanner *bufio.Scanner, separator rune) (CSVData, error) {
	var arr CSVData

	for scanner.Scan() {
		line := scanner.Text()

		if separator == 0 {
			separator = detectSeparator(line)
		}

		values := strings.Split(line, string(separator))

		interfaceValues := make([]interface{}, len(values))
		for i, v := range values {
			interfaceValues[i] = inferType(v)
		}

		arr = append(arr, interfaceValues)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return arr, nil
}

// Parse parses CSV data from a byte slice
func Parse(csv []byte, separator rune) (CSVData, error) {
	scanner := bufio.NewScanner(bytes.NewReader(csv))
	return parseScanner(scanner, separator)
}

// ParseFile parses CSV data from a file
func ParseFile(path string, separator rune) (CSVData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return parseScanner(scanner, separator)
}

// ParseCSVToStruct parses CSV data from a byte slice and unmarshals it into the provided struct type
func ParseCSVToStruct(csvData []byte, separator rune, result interface{}) error {
	data, err := Parse(csvData, separator)
	if err != nil {
		return err
	}

	jsonData, err := ConvertToJSON(data, true)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, result)
}

// ParseFileToStruct parses CSV data from a file and unmarshals it into the provided struct type
func ParseFileToStruct(filePath string, separator rune, result interface{}) error {
	data, err := ParseFile(filePath, separator)
	if err != nil {
		return err
	}

	jsonData, err := ConvertToJSON(data, true)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, result)
}

// ConvertToJSON converts CSVData to JSON format, using headers if specified
func ConvertToJSON(data CSVData, header bool) ([]byte, error) {
	if header && len(data) > 0 {
		headers := data[0]
		// Remove trailing nil values from headers, if any
		for i := len(headers) - 1; i >= 0; i-- {
			if headers[i] == nil {
				headers = headers[:i]
			} else {
				break
			}
		}

		var result []map[string]interface{}

		for _, row := range data[1:] {
			rowMap := make(map[string]interface{})
			for i, value := range row[:len(headers)] {
				key, ok := headers[i].(string)
				if ok {
					rowMap[key] = value
				}
			}
			result = append(result, rowMap)
		}

		return json.Marshal(result)
	}

	return json.Marshal(data)
}

// detectSeparator detects the separator used in a CSV line, pretty basic
func detectSeparator(line string) (separator rune) {
	separators := []rune{',', ';', '\t', '|'}
	maxCount := 0

	for _, sep := range separators {
		count := strings.Count(line, string(sep))
		if count > maxCount {
			maxCount = count
			separator = sep
		}
	}

	if separator == 0 {
		separator = ','
	}

	return
}

func inferType(value string) interface{} {
	// Return nil if value is "null"
	if value == "null" || value == "" {
		return nil
	}

	// Attempt to parse as bool
	if inferredValue, err := strconv.ParseBool(value); err == nil {
		return inferredValue
	}

	// Attempt to parse as int
	if inferredValue, err := strconv.Atoi(value); err == nil {
		return inferredValue
	}

	// Replace commas with dots for float parsing
	valueFloat := strings.Replace(value, ",", ".", 1)
	// Attempt to parse as float64
	if inferredValue, err := strconv.ParseFloat(valueFloat, 64); err == nil {
		return inferredValue
	}

	return value
}
