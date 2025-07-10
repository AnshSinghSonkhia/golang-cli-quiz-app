// Quiz CLI App: Reads math problems from a CSV, randomizes, and quizzes the user with a timer.
package main

// Import required packages
import (
	"encoding/csv" // for reading CSV files
	"flag"         // for command-line flags
	"fmt"          // for formatted I/O
	"math/rand"    // for randomizing questions
	"os"           // for file and exit operations
	"time"         // for timer and seeding randomness
)

// problemPuller reads all problems from the given CSV file and returns a slice of problem structs.
func problemPuller(fileName string) ([]problem, error) {
	// Open the CSV file
	fObj, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	csvReader := csv.NewReader(fObj)
	cLines, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading csv file: %w", err)
	}

	// Parse the CSV lines into problem structs
	return parseProblem(cLines), nil
}

// parseProblem converts CSV lines to a slice of problem structs.
func parseProblem(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i, line := range lines {
		if len(line) != 2 {
			continue // skip lines that do not have exactly two elements
		}
		r[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}
	return r
}

// problem represents a single quiz question and its answer.
type problem struct {
	question string // The question to ask (e.g., "2+2")
	answer   string // The correct answer (e.g., "4")
}

// main is the entry point of the quiz app.
func main() {
	// Parse command-line flags for CSV file and timer
	fileName := flag.String("f", "quiz.csv", "path of csv file")
	timer := flag.Int("t", 30, "time limit in seconds")
	flag.Parse()

	// Read problems from the CSV file
	problems, err := problemPuller(*fileName)
	if err != nil {
		exit(fmt.Sprintf("Error reading the file: %s\n", err.Error()))
	}

	// Shuffle the problems randomly
	randomizeProblems(problems)

	// Select only 5 random questions (or fewer if not enough problems)
	questionsToAsk := 5
	if len(problems) < questionsToAsk {
		questionsToAsk = len(problems)
	}
	selectedProblems := problems[:questionsToAsk]

	correctCount := 0                                          // Track number of correct answers
	tObj := time.NewTimer(time.Duration(*timer) * time.Second) // Quiz timer
	answerCh := make(chan string)                              // Channel for user answers

problemLoop:
	for i, p := range selectedProblems {
		var answer string
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		// Read user input in a goroutine so we can select on timer
		go func() {
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop // Exit if time is up
		case answer = <-answerCh:
			if answer == p.answer {
				correctCount++ // Increment if correct
			} else {
				fmt.Printf("Wrong! The correct answer is %s\n", p.answer)
			}

			// Close channel if last problem
			if i == len(selectedProblems)-1 {
				close(answerCh)
			}
		}
	}

	// Print final score and wait for user to press Enter
	fmt.Printf("You scored %d out of %d\n", correctCount, questionsToAsk)
	fmt.Printf("Press Enter to exit...\n")
	<-answerCh // Wait for Enter before exiting
}

// randomizeProblems shuffles the problems slice in place using Fisher-Yates algorithm.
func randomizeProblems(problems []problem) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range problems {
		j := r.Intn(len(problems))
		problems[i], problems[j] = problems[j], problems[i]
	}
}

// exit prints the error message and exits the program.
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
