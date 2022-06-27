package main

import (
  "encoding/csv"
	"flag"
  "fmt"
  "io"
  "log"
  "os"
  "time"
)

type problem struct {
  q string
  a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file with the format of question,answer")
	timeLimit := flag.Int("time", 15, "the time limit for the quiz in seconds")
	flag.Parse()

  problems := readCSVFile(*csvFilename)
  timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

  numCorrect, numQuestions := 0, len(problems)

problemloop:
  for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i + 1, problem.q)

		answerCh := make(chan string)

		go func() {
			var str string
			fmt.Scanf("%s", &str)
			answerCh <- str
		}()

		select {
		case answer := <-answerCh:
			if answer == problem.a {
				numCorrect++
			}
		case <-timer.C:
			fmt.Println()
			break problemloop
		}
	}
			
	fmt.Printf("You scored %d out of %d.\n", numCorrect, numQuestions)
}

func readCSVFile(filePath string) []problem {
  problems := make([]problem, 0)

  file, err := os.Open(filePath)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  
  csvReader := csv.NewReader(file)
  for {
    line, err := csvReader.Read()
    
    if err == io.EOF {
      break
    }

    if err != nil {
      log.Fatal(err)
    }
    
    problems = append(problems, problem{line[0], line[1]})
  }

  return problems
}

