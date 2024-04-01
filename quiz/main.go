package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"log"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "file to get questions from")
	timer := flag.Int("limit", 30, "the time limit in seconds for thr quiz")
	shuffle := flag.Bool("shuffle", false, "Whether the question should be shuffled")
	flag.Parse()

	csvFile, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatalf("can't open the file %s: %v", *csvFileName, err)

	}
	defer csvFile.Close()

	csvEncoder := csv.NewReader(csvFile)
	lines, err := csvEncoder.ReadAll()
	if err != nil {
		return log.Fatalf("can't parse csv: %v", err)
	}

	problems := parseLines(lines)
	err = quiz(problems, *timer, *shuffle)
	if err != nil {
		log.Fatalf("quiz: %v", err)
	}
}

func quiz(problems []problem, timeLimit int, shuffle bool) error {

	if shuffle {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Sscanln("Press enter to start a quiz")

	var correctAnswers, questionsAsked int
	answers := make(chan string, 0)
	timer := time.NewTimer(time.Second * time.Duration(timeLimit))

	for _, problem := range problems {
		fmt.Printf("%s = ", problem.question)
		questionsAsked++

		go func() {
			var answer string
			fmt.Scanln(&answer)
			answers <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime is out")
			fmt.Printf("Result: %d out of %d", correctAnswers, questionsAsked)
			return nil
		case answer := <-answers:
			if strings.TrimSpace(answer) == problem.answer {
				correctAnswers++
			}
		}
	}

	fmt.Printf("Result: %d out of %d", correctAnswers, questionsAsked)
	return nil
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i].question = line[0]
		problems[i].answer = strings.TrimSpace(line[1])
	}
	return problems
}
