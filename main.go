package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func main() {
	// Define a flag for the input file
	inputFile := flag.String("input", "", "Input markdown file")
	flag.Parse()

	// Check if the input file is provided
	if *inputFile == "" {
		log.Fatal("Please provide an input file using the -input flag")
	}

	// Read the input file
	md, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Convert markdown to HTML
	html := mdToHTML(md)

	// Write the HTML to output.html
	err = ioutil.WriteFile("output.html", html, 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Conversion complete. Output written to output.html")
}
