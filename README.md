# Golang CLI Quiz App

![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![CSV](https://img.shields.io/badge/CSV-Data-green)
![CLI](https://img.shields.io/badge/CLI-Command%20Line-blue)
![Concurrency](https://img.shields.io/badge/Concurrency-Goroutines-yellow)
![Channels](https://img.shields.io/badge/Channels-Go%20Channels-orange)

A simple command-line quiz application written in Go. This app reads math problems from a CSV file, randomizes the questions, and quizzes the user with a configurable timer. It's a great project for learning Go basics, file I/O, concurrency, and command-line flag parsing.

## Features
- Reads quiz questions and answers from a CSV file
- Randomizes the order of questions each run
- Asks only a subset (default: 5) of questions per session
- Enforces a time limit for the entire quiz (default: 30 seconds)
- Provides instant feedback on each answer
- Displays the final score at the end

## How It Works
1. The app loads questions from a CSV file (default: `quiz.csv`).
2. It shuffles the questions and selects 5 random ones.
3. The user is prompted to answer each question, one at a time.
4. The quiz ends when either all questions are answered or the timer runs out.
5. The final score is displayed.

## CSV Format
The CSV file should have two columns per line:
```
question,answer
70+20,90
15+25,40
...etc
```
No header row is required. Each line represents one question and its correct answer.

## Usage
### Build
```sh
go build -o quizapp main.go
```

### Run
```sh
./quizapp -f quiz.csv -t 30
```
- `-f` : Path to the CSV file (default: `quiz.csv`)
- `-t` : Time limit for the quiz in seconds (default: 30)

### Example
```
Problem #1: 15+25 = 40
Problem #2: 70+20 = 90
Wrong! The correct answer is 90
...
You scored 3 out of 5
Press Enter to exit...
```

## Project Structure
- `main.go` : Main application source code
- `quiz.csv` : Sample CSV file with quiz questions
- `README.md` : Project documentation (this file)

## How to Add More Questions
Just add more lines to `quiz.csv` in the format:
```
question,answer
```

## Customization
- Change the number of questions asked by modifying the `questionsToAsk` variable in `main.go`.
- Adjust the timer using the `-t` flag.

## Learning Highlights
- File I/O with Go
- CSV parsing
- Command-line flag parsing
- Goroutines and channels for concurrency
- Randomization and shuffling
- Error handling and user feedback

## License
This project is open source and free to use for learning and personal projects.

## Use of Goroutines

### Why we used goroutines
Goroutines are used to allow the program to wait for user input and the timer at the same time. Without goroutines, the program would block on user input and never notice if the timer expired.

### Where we used goroutines
Inside the main quiz loop, for each question, we use a goroutine to read the user's answer:

```go
    go func() {
        fmt.Scanf("%s\n", &answer)
        answerCh <- answer
    }()
```

### How we used goroutines
- We launch a goroutine that waits for the user's input and sends it to the `answerCh` channel.
- Meanwhile, the main goroutine uses a `select` statement to wait for either:
  - The timer to expire (`<-tObj.C`), or
  - The user's answer to arrive on the channel (`answer = <-answerCh`)
- This allows the quiz to end immediately if the timer runs out, even if the user hasn't answered yet.

**Summary:**
Goroutines enable concurrent waiting for both user input and the timer, making the quiz responsive and time-limited.

## Use of Channels

### Why we used channels
Channels are used to safely communicate between goroutines. In this app, a channel is needed to send the user's answer from the goroutine (that reads input) back to the main goroutine, which is waiting for either the answer or the timer to expire.

### Where we used channels
We declare and use the channel in the main quiz loop:

```go
answerCh := make(chan string) // Channel for user answers
```
We send the answer into the channel from the goroutine:
```go
answerCh <- answer
```
And we receive the answer in the main goroutine using a select statement:
```go
case answer = <-answerCh:
```

### How we used channels
- The channel `answerCh` is created before the quiz loop starts.
- For each question, a goroutine reads user input and sends it to `answerCh`.
- The main goroutine waits for either a value from `answerCh` (user answered) or the timer to expire.
- This pattern ensures safe and synchronized communication between concurrent parts of the program.

**Summary:**
Channels allow the main quiz logic and the user input goroutine to communicate efficiently and safely, enabling the quiz to be both interactive and time-limited.
