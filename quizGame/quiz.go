package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exitStr := fmt.Sprintf("failed to open the CSV file %s\n", *csvFilename)
		exit(exitStr)
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()

	if err != nil {
		exitStr := fmt.Sprintf("failed to read the CSV file %s\n", *csvFilename)
		exit(exitStr)
	}

	problems := parseLines(lines)

	correct := 0

	// try breaking this into its own function and writing a test!
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		var answer string

		fmt.Scanf("%s\n", &answer)

		if answer == p.a {
			fmt.Println("Correct!")
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func exit(msg string) {
	fmt.Println(msg)

	os.Exit(1)
}

// better to depend on this type rather than the 2D slice we get back from parsing CSV.  This way,
// in the future our program could parse different kinds of input, and the rest of the code would work fine
// provided we manipulate the data into this shape first.
type problem struct {
	q string
	a string
}
