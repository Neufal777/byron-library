package main

import "github.com/Byron/core"

type Yo struct {
	Name string
	Age  string
}

func main() {

	Categories := []string{
		//"machine learning",
		//"artificial Inteligence",
		"neurociencia",
		//"inteligencia artificial",
		//"math",
		//"maths",
		//"biology",
		//"biologia",
	}

	for _, c := range Categories {
		core.LIBGENDownloadAll(c)

	}

}
