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
        :root {
            --bg-color: #ffffff;
            --text-color: #333333;
            --link-color: #0366d6;
            --code-bg: #f6f8fa;
            --blockquote-color: #6a737d;
            --border-color: #eaecef;
            --button-bg: #f1f1f1;
            --button-color: #333333;
        }

        body.dark-mode {
            --bg-color: #1e1e1e;
            --text-color: #d4d4d4;
            --link-color: #4da6ff;
            --code-bg: #2d2d2d;
            --blockquote-color: #a0a0a0;
            --border-color: #4a4a4a;
            --button-bg: #3a3a3a;
            --button-color: #d4d4d4;
        }

        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: var(--text-color);
            background-color: var(--bg-color);
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            transition: background-color 0.3s ease, color 0.3s ease;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }
        h1 { font-size: 2em; border-bottom: 1px solid var(--border-color); padding-bottom: .3em; }
        h2 { font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: .3em; }
        a { color: var(--link-color); text-decoration: none; }
        a:hover { text-decoration: underline; }
        pre {
            background-color: var(--code-bg);
            border-radius: 3px;
            font-size: 85%;
            line-height: 1.45;
            overflow: auto;
            padding: 16px;
            position: relative;
        }
        code {
            background-color: var(--code-bg);
            border-radius: 3px;
            font-size: 85%;
            margin: 0;
            padding: .2em .4em;
        }
        blockquote {
            border-left: .25em solid var(--border-color);
            color: var(--blockquote-color);
            padding: 0 1em;
        }
        img { max-width: 100%; }
        #mode-toggle {
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 10px;
            background-color: var(--button-bg);
            color: var(--button-color);
            border: 1px solid var(--border-color);
            border-radius: 5px;
            cursor: pointer;
        }
        .copy-button {
            position: absolute;
            top: 5px;
            right: 5px;
            padding: 5px 10px;
            background-color: var(--button-bg);
            color: var(--button-color);
            border: none;
            border-radius: 3px;
            cursor: pointer;
        }
        .copy-button:hover {
            opacity: 0.8;
        }
    </style>
</head>
<body>
    <button id="mode-toggle">Toggle Dark Mode</button>
    {{.Body}}
    <script>
        // Dark mode toggle functionality
        const modeToggle = document.getElementById('mode-toggle');
        const body = document.body;
        const currentTheme = localStorage.getItem('theme');

        if (currentTheme) {
            body.classList.add(currentTheme);
        }

        modeToggle.addEventListener('click', () => {
            body.classList.toggle('dark-mode');
            const theme = body.classList.contains('dark-mode') ? 'dark-mode' : '';
            localStorage.setItem('theme', theme);
        });

        // Code copy functionality
        document.querySelectorAll('pre').forEach((block) => {
            const button = document.createElement('button');
            button.innerText = 'Copy';
            button.className = 'copy-button';
            button.addEventListener('click', () => {
                const code = block.querySelector('code');
                navigator.clipboard.writeText(code.innerText).then(() => {
                    button.innerText = 'Copied!';
                    setTimeout(() => {
                        button.innerText = 'Copy';
                    }, 2000);
                }, (err) => {
                    console.error('Failed to copy: ', err);
                });
            });
            block.appendChild(button);
        });
    </script>
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
