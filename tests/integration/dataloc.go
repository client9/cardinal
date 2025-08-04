// Package dataloc provides functionality to find the source code location of
// table-driven test cases.
package integration

import (
	"bytes"
	"fmt"
	//	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// L returns the source code location of the test case identified by its name.
// It attempts runtime source code analysis to find the location
// by using the expression passed to dataloc.L().
// So some restrictions apply:
//   - The function must be invoked as "dataloc.L".
//   - The argument must be an expression of the form "dataloc.L(testcase.key)"
//     , where "testcase" is a variable declared as "for _, testcase := range testcases"
//     , and "testcases" is a slice of a struct type
//     , whose "key" field is a string which is passsed to L().
//   - or "dataloc.L(key)"
//     , where key is a variable declared as "for key, value := range testcases"
//     , and "testcases" is a map of string to any type
//     , and "key" is the string which is passed to L().
//
// See Example.
func L(name string) string {
	return loc(name, 2)
}

func L3(name string) string {
	return loc(name, 3)
}

func L4(name string) string {
	return loc(name, 4)
}

func L5(name string) string {
	return loc(name, 5)
}

func L6(name string) string {
	return loc(name, 6)
}

func loc(value string, step int) string {
	_, file, _, _ := runtime.Caller(step)
	//	log.Printf("Caller Step %d: %s %d", step, file, line)
	cwd, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	file, err = filepath.Rel(cwd, file)
	if err != nil {
		return err.Error()
	}

	token := []byte(strconv.Quote(value))
	fbytes, err := os.ReadFile(file)
	num := 0
	for line := range bytes.SplitSeq(fbytes, []byte{'\n'}) {
		num++
		if bytes.Contains(line, token) {
			return fmt.Sprintf("%s:%d", file, num)
		}
	}
	return "(unknown)"
}
