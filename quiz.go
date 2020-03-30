// Read csv file, print out contents in loop

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type options struct {
	shuffle bool
	foo     bool
}

func main() {
	options := handleOptions()
	questions := getQuestions("problems.csv")
	if options.shuffle {
		questions = shuffle(questions)
	}

	updatedQuestions := askQuestions(questions)
	fmt.Printf("%+v\n", correctAnswers(updatedQuestions))
	printScore(questions, correctAnswers(updatedQuestions))
}

func handleOptions() options {
	shuffle := flag.Bool("shuffle", false, "Suffle questions")
	foo := flag.Bool("foo", false, "foo")
	flag.Parse()

	return options{shuffle: *shuffle, foo: *foo}
}

func printScore(questionsAsked []QandA, questionsCorrect []QandA) {
	fmt.Printf("You scored %v out of %v\n", len(questionsCorrect), len(questionsAsked))
}

func askQuestions(questions []QandA) []QandA {
	var updatedQuestions []QandA

	timeout := time.After(3 * time.Second)
adam:
	for _, question := range questions {
		done := make(chan bool)

		go func(question *QandA) {
			val := ask(*question)
			question.setAnswer(val)
			done <- true
		}(&question)

		select {
		case <-done:
			fmt.Println(question)
		case <-timeout:
			fmt.Println("Sorry, out of time")
			break adam
		}
		updatedQuestions = append(updatedQuestions, question)
	}
	return updatedQuestions

}

func ask(question QandA) string {
	fmt.Print(question.question + ": ")
	var input string

	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal(err)
	}

	return input
}

func getQuestions(csvFileName string) []QandA {
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(csvFile)

	// Since we know the leng
	//
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

func shuffle(a []QandA) []QandA {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}
