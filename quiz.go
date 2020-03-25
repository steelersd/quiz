// Read csv file, print out contents in loop

package main

import (
	// "flag"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	questions := getQuestions("problems.csv")

	updatedQuestions := askQuestions(questions)
	fmt.Printf("%+v\n", correctAnswers(updatedQuestions))
	printScore(questions, correctAnswers(updatedQuestions))
}

func printScore(questionsAsked []QandA, questionsCorrect []QandA) {
	fmt.Printf("You scored %v out of %v\n", len(questionsCorrect), len(questionsAsked))
}

func askQuestions(questions []QandA) []QandA {
	var updatedQuestions []QandA

	for _, question := range questions {
		fmt.Print(question.question + ": ")
		var input string

		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		question.setAnswer(input)
		updatedQuestions = append(updatedQuestions, question)
		// fmt.Printf("question: %v answer: %v, %v \n", question.question, question.answer, question.correct)
	}
	return updatedQuestions
}

func getQuestions(csvFileName string) []QandA {
	csvFile, _ := os.Open(csvFileName)
	r := csv.NewReader(csvFile)

	var questions []QandA

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		questions = append(questions, QandA{question: record[0], answer: record[1]})
	}
	return questions
}

// QandA struct
type QandA struct {
	question  string
	answer    string
	correct   bool
	userInput interface{}
}

func (qa *QandA) setAnswer(answer string) {
	qa.userInput = answer
	val, err := strconv.Atoi(qa.answer)
	if err != nil {
		log.Fatal(err)
	}
	val2, e := strconv.Atoi(answer)
	if e != nil {
		log.Fatal(e)
	}
	qa.correct = val2 == val
}

func correctAnswers(qas []QandA) []QandA {
	filtered := filter(qas, func(qa QandA) bool {
		return qa.correct
	})
	return filtered
}

func filter(vs []QandA, f func(QandA) bool) []QandA {
	vsf := make([]QandA, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
