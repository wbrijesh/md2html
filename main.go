package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown Converted to HTML</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }
        h1 { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: .3em; }
        h2 { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: .3em; }
        a { color: #0366d6; text-decoration: none; }
        a:hover { text-decoration: underline; }
        pre {
            background-color: #f6f8fa;
            border-radius: 3px;
            font-size: 85%;
            line-height: 1.45;
            overflow: auto;
            padding: 16px;
        }
        code {
            background-color: rgba(27,31,35,.05);
            border-radius: 3px;
            font-size: 85%;
            margin: 0;
            padding: .2em .4em;
        }
        blockquote {
            border-left: .25em solid #dfe2e5;
            color: #6a737d;
            padding: 0 1em;
        }
        img { max-width: 100%; }
    </style>
</head>
<body>
    {{.Body}}
</body>
</html>
`

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
	md, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Convert markdown to HTML
	htmlBody := mdToHTML(md)

	// Create a template
	tmpl, err := template.New("webpage").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing HTML template: %v", err)
	}

	// Create the output file
	outputFile, err := os.Create("output.html")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Execute the template
	err = tmpl.Execute(outputFile, struct{ Body template.HTML }{template.HTML(htmlBody)})
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Println("Conversion complete. Output written to output.html")
}
