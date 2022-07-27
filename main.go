package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// plays a game based on a file containing the questions
// and answers. the questions can range from general knowledge
// to arithmetic, the only condition is that the answer must
// match and the final score is calculated.

// the program takes two input flags: -f and -l which take in
// the filename containing the q & a, and the time limit in seconds
// to finish the quiz.
func main() {
	var fileName string
	var limit int
	flag.StringVar(&fileName, "f", "challenge.csv", "name of file containing the challenge")
	flag.IntVar(&limit, "l", 60, "time limit of the challenge")
	flag.Parse()
	reader := csv.NewReader(readFile(fileName))
	data, err := reader.ReadAll()
	handleErr(err)
	score := 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(limit))
	defer cancel()
	score = play(ctx, data)
	fmt.Printf("You Scored: %v / %v\n", score, len(data))
}

func play(ctx context.Context, data [][]string) int {
	score := 0
	if ctx != nil {
		if ctx.Err() != nil {
			return score
		}
	}
	go func() {
		select {
		case <-ctx.Done():
			fmt.Printf("You Scored: %v / %v\n", score, len(data))
			os.Exit(1)
		}
	}()
	for _, challenge := range data {
		question, ans := challenge[0], challenge[len(challenge)-1]
		var userInput string
		fmt.Printf("Question: %v, Answer: ", question)
		fmt.Scan(&userInput)
		if strings.EqualFold(ans, userInput) {
			score += 1
		}
	}
	return score
}

func readFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	handleErr(err)
	return file
}

func handleErr(err error) {
	if err != nil {
		os.Exit(1)
	}
}
