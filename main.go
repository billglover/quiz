package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("file", "problems.csv", "CSV file containing questions and answers")
	duration := flag.Int("duration", 30, "duration of the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "randomises the order of questions")
	flag.Parse()

	err := runQuiz(*filename, *duration, *shuffle)
	if err != nil {
		fmt.Println("failed to run quiz:", err)
		os.Exit(1)
	}
}

func runQuiz(fn string, d int, shuffle bool) error {
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

	if shuffle {
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	score := 0
	count := len(records)

	fmt.Printf("Quiz will last %d seconds.\n", d)
	fmt.Print("Press Enter to begin...")
	fmt.Scanln()
	fmt.Println()

	timerChan := time.After(time.Second * time.Duration(d))
	inputChan := make(chan string)

	for _, l := range records {

		fmt.Printf("Q: %s\n", l[0])
		go getAnswer(inputChan)

		select {

		case <-timerChan:
			fmt.Printf("\nTime's up!\n")
			displayScore(score, count)
			return nil

		case answer := <-inputChan:
			fmt.Println()
			if strings.ToLower(answer) != strings.ToLower(l[1]) {
				continue
			}
			score++
		}
	}

	displayScore(score, count)
	return nil
}

func displayScore(score, count int) {
	fmt.Println("---")
	fmt.Printf("You scored %d out of %d!\n", score, count)
}

func getAnswer(returnChan chan<- string) {
	fmt.Print("A: ")

	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		return
	}

	returnChan <- strings.TrimSpace(answer)
}
