package repl_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const slash = string(os.PathSeparator)

// Runs make_env and then compares output.
func TestMakeEnv(t *testing.T) {

	got, err := exec.Command("go", "build", "-o", "make_env", "make_env.go").Output()
	if err != nil {
		fmt.Printf("Error running: go build -o make_env make_env.go: %s %s", got, err)
		log.Fatal(err)
	}

	got, err = exec.Command("./make_env").Output()
	if err != nil {
		fmt.Printf("Error running ./make_env: %s", got)
		log.Fatal(err)
	}


	rightFileName := "eval_imports.go"

	var rightFile *os.File
	rightFile, err = os.Open(rightFileName) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 500000)
	count, err := rightFile.Read(data)
	if err != nil {
		t.Errorf("Failed to read 'right' data file %s:", rightFileName)
		log.Fatal(err)
	}
	want := string(data[0:count])
	if string(got) != want {
		gotName := fmt.Sprintf("testdata%snew_eval_imports.go",  slash)
		gotLines  := strings.Split(string(got), "\n")
		wantLines := strings.Split(string(want), "\n")
		wantLen   := len(wantLines)
		for i, line := range(gotLines) {
			if i == wantLen {
				fmt.Println("want results are shorter than got results, line", i+1)
				break
			}
			if line != wantLines[i] {
				fmt.Println("results differ starting at line", i+1)
				fmt.Println("got:\n", line)
				fmt.Println("want:\n", wantLines[i])
				break
			}
		}
		if err := ioutil.WriteFile(gotName, got, 0666); err == nil {
			fmt.Printf("Full results are in file %s\n", gotName)
		}
		t.Errorf("make_env comparison test failed")
	}

	// Print a helpful hint if we don't make it to the end.
	hint := "Run manually"
	defer func() {
		if hint != "" {
			fmt.Println("FAIL")
			fmt.Println(hint)
		} else {
			fmt.Println("PASS")
		}
	}()

	hint = "" // call off the hounds
	return
}
