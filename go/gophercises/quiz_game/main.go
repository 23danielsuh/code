package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func hello() {
	fmt.Println("Hello world!")
}

func readCSVFile(filePath string) {
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
    
    fmt.Println(line[0])
  }
}

func main() {
  readCSVFile("problems.csv")
}
