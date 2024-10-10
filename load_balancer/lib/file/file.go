package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reads the first line from a file and removes it from the file
func ReadFirstLineAndRemove(filename string) (string, error) {
	return readLineAndRemove(filename, true)
}

// Reads the last line from a file and removes it from the file
func ReadLastLineAndRemove(filename string) (string, error) {
	return readLineAndRemove(filename, false)
}

// Reads either the first or last line from a file and removes it
func readLineAndRemove(filename string, first bool) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if len(lines) == 0 {
		return "", nil
	}

	var line string
	if first {
		line = lines[0]
		lines = lines[1:]
	} else {
		line = lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	}

	// Truncate the file
	if err := file.Truncate(0); err != nil {
		return "", err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Write the remaining lines back to the file
	writer := bufio.NewWriter(file)
	for _, l := range lines {
		fmt.Fprintln(writer, l)
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// Appends a line to a file
func AppendToFile(filename, line string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = fmt.Fprintln(writer, line)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// readIPAddresses reads IP addresses from the given file
func ReadIPAddresses(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ipAddresses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipAddresses = append(ipAddresses, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ipAddresses, nil
}
