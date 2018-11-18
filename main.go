package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	filename := flag.String("file", "problems.csv", "CSV file containing questions and answers")
	flag.Parse()

	err := runQuiz(*filename)
	if err != nil {
		fmt.Println("failed to run quiz:", err)
		os.Exit(1)
	}
}

func runQuiz(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	score := 0
	count := len(records)

	for _, l := range records {

		fmt.Printf("Q: %s\n", l[0])

		var answer string
		fmt.Print("A: ")
		_, err := fmt.Scanln(&answer)
		fmt.Println()
		if err != nil {
			// NOTE: swallowing errors is usually bad, but we
			// choose not to display the error as it is for the
			// user to figure out why their answers are incorrect.
			continue
		}

		if answer != l[1] {
			continue
		}

		score++
	}

	fmt.Println("---")
	fmt.Printf("You scored %d out of %d!\n", score, count)

	return nil
}
