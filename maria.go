// Maria is a simulation of the Marie machine described in chapter 4 of
// "Computer Organization and Architecture" by Linda Null and Julia Lobur.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: maria file")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	m := &Machine{}
	m.Load(f)
	m.Run()
}
