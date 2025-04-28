package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	qs := initQState()
	fileFlag := defineFlags()
	flag.Parse()
	fmt.Printf("using file: %s\n", *fileFlag)

	fileName := *fileFlag

	file, err := os.Open(fileName)
	logFatalIfErr(err)
	// open file, if error opening file just log the error and exit
	fmt.Println("File successfully opened")

	data, err := getData(file)
	logFatalIfErr(err)
	//parse information out from csv file

	scanner := bufio.NewScanner(os.Stdin)
	for _, row := range data {
		qna := mapQna(row[0], row[1])
		runQuestion(scanner, qna, qs)
		qs.printCurrState()
	}
	defer file.Close()
}

type qna struct {
	question string
	answer   string
}

type quizState struct {
	questionAttempted int
	numCorrect        int
}

func runQuestion(scanner *bufio.Scanner, qna *qna, qs *quizState) {
	qna.printQuestion()
	for {
		if scanner.Scan() {
			input := scanner.Text()
			switch input == qna.getAnswer() {
			case true:
				fmt.Println("Correct!")
				qs.incNumCorrect()
			case false:
				fmt.Println("Wrong!")
			}
			qs.incQAttempted()
			break
		}
	}
}

func initQState() *quizState {
	state := quizState{
		0, 0,
	}
	return &state
}

func (qs *quizState) incQAttempted() {
	qs.questionAttempted++
}
func (qs *quizState) incNumCorrect() {
	qs.numCorrect++
}

func (qs *quizState) printCurrState() {
	fmt.Println("questions attempted: ", qs.questionAttempted)
	fmt.Println("number of correct answers:", qs.numCorrect)
}

func mapQna(question, answer string) *qna {
	currentQna := qna{
		question,
		answer,
	}
	return &currentQna
}

func (qna *qna) printQuestion() {
	question := fmt.Sprintf("What is %s?", qna.getQuestion())
	fmt.Println(question)
}

func (q *qna) getQuestion() string {
	return q.question
}

func (q *qna) getAnswer() string {
	return q.answer
}

func logFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func defineFlags() (filename *string) {
	const (
		defaultFile = "problems.csv"
		usage       = "the filename where the quiz details are"
	)
	fileFlag := flag.String("quiz_file", defaultFile, usage)
	return fileFlag
}

func getData(file *os.File) (data [][]string, err error) {
	reader := csv.NewReader(file)
	// note fhat file fulfils the interface io.Reader, but the reader that is returned is an instance of *csv.Reader
	reader.FieldsPerRecord = 2
	data, err = reader.ReadAll()
	if err != nil {
		return nil, err
	}
	// set to variable number of fields by declaring negative value for FieldsPerRecord
	if len(data) == 0 {
		return nil, errors.New("no data was returned from reading file")
	}
	return data, nil
}
