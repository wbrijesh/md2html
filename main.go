package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const css = `
:root {
    --bg-color: #ffffff;
    --text-color: #24292e;
    --link-color: #0366d6;
    --code-bg: #f6f8fa;
    --blockquote-color: #6a737d;
    --blockquote-border: #dfe2e5;
    --table-border: #dfe2e5;
    --table-bg: #f6f8fa;
    --button-border: #d1d5da;
}

body.dark-mode {
    --bg-color: #1e1e1e;
    --text-color: #d4d4d4;
    --link-color: #4da6ff;
    --code-bg: #2d2d2d;
    --blockquote-color: #a0a0a0;
    --blockquote-border: #4a4a4a;
    --table-border: #4a4a4a;
    --table-bg: #2d2d2d;
    --button-border: #4a4a4a;
}

body {
    font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji";
    font-size: 16px;
    line-height: 1.5;
    word-wrap: break-word;
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    color: var(--text-color);
    background-color: var(--bg-color);
    transition: background-color 0.3s ease, color 0.3s ease;
}

h1, h2, h3, h4, h5, h6 {
    margin-top: 24px;
    margin-bottom: 16px;
    font-weight: 600;
    line-height: 1.25;
}

h1 { font-size: 2em; }
h2 { font-size: 1.5em; }
h3 { font-size: 1.25em; }
h4 { font-size: 1em; }
h5 { font-size: 0.875em; }
h6 { font-size: 0.85em; }

a {
    color: var(--link-color);
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

pre, code {
    background-color: var(--code-bg);
    border-radius: 3px;
    font-size: 16px;
}

pre {
    line-height: 1.45;
    overflow: auto;
    padding: 16px;
}

code {
    padding: .2em .4em;
}

pre code {
    background-color: transparent;
    padding: 0;
}

blockquote {
    border-left: .25em solid var(--blockquote-border);
    color: var(--blockquote-color);
    padding: 0 1em;
}

table {
    border-collapse: collapse;
    width: 100%;
}

table th, table td {
    border: 1px solid var(--table-border);
    padding: 6px 13px;
}

table tr:nth-child(2n) {
    background-color: var(--table-bg);
}

img {
    max-width: 100%;
    box-sizing: content-box;
}

.mode-toggle {
    position: fixed;
    top: 10px;
    right: 10px;
    padding: 5px 10px;
    background-color: var(--bg-color);
    color: var(--text-color);
    border: 1px solid var(--button-border);
    border-radius: 5px;
    cursor: pointer;
    transition: border-color 0.3s ease;
}

.copy-button {
    position: absolute;
    top: 5px;
    right: 5px;
    padding: 5px 10px;
    background-color: var(--code-bg);
    border: 1px solid var(--button-border);
    border-radius: 3px;
    cursor: pointer;
    color: var(--text-color);
    transition: border-color 0.3s ease;
}

.copy-button:hover {
    opacity: 0.8;
}

pre {
    position: relative;
}
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

func convertToHtmlWithFeatures(filename string) (string, error) {
	md, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	html := mdToHTML(md)

	outputFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".html"

	// JavaScript for copy functionality and dark mode toggle
	script := `
    <script>
    document.addEventListener('DOMContentLoaded', (event) => {
        // Copy button functionality
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

        // Dark mode toggle functionality
        const toggleButton = document.getElementById('mode-toggle');
        const body = document.body;
        const currentTheme = localStorage.getItem('theme');

        if (currentTheme) {
            body.classList.add(currentTheme);
        }

        toggleButton.addEventListener('click', () => {
            body.classList.toggle('dark-mode');
            const theme = body.classList.contains('dark-mode') ? 'dark-mode' : '';
            localStorage.setItem('theme', theme);
            updatePrismTheme(theme);
        });

        function updatePrismTheme(theme) {
            const prismTheme = document.getElementById('prism-theme');
            if (theme === 'dark-mode') {
                prismTheme.href = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/themes/prism-okaidia.min.css';
            } else {
                prismTheme.href = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/themes/prism.min.css';
            }
        }

        updatePrismTheme(currentTheme);
    });
    </script>
    `

	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    <style>
%s
    </style>
    <link id="prism-theme" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/themes/prism.min.css" rel="stylesheet" />
</head>
<body>
    <button id="mode-toggle" class="mode-toggle">Toggle Dark Mode</button>
%s
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/components/prism-core.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/plugins/autoloader/prism-autoloader.min.js"></script>
%s
</body>
</html>`, filename, css, string(html), script)

	err = os.WriteFile(outputFilename, []byte(fullHTML), 0644)
	if err != nil {
		return "", err
	}

	fmt.Printf("Converted %s to %s\n", filename, outputFilename)
	return outputFilename, nil
}

func selectMarkdownFile(files []string) (string, error) {
	var selectedFile string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Markdown file to convert").
				Options(huh.NewOptions(files...)...).
				Value(&selectedFile),
		),
	)

	err := form.Run()
	if err != nil {
		return "", fmt.Errorf("error in form: %v", err)
	}

	if selectedFile == "" {
		return "", fmt.Errorf("no file selected")
	}

	return selectedFile, nil
}

func confirmOpenFile() (bool, error) {
	var openFile bool
	openForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to open the converted HTML file?").
				Value(&openFile),
		),
	)

	err := openForm.Run()
	if err != nil {
		return false, fmt.Errorf("error in open file form: %v", err)
	}

	return openFile, nil
}

func openFileInBrowser(filename string) error {
	cmd := exec.Command("open", filename)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	fmt.Printf("Opened %s in your default browser.\n", filename)
	return nil
}

func listMarkdownFiles() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var mdFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			mdFiles = append(mdFiles, file.Name())
		}
	}
	return mdFiles, nil
}

func main() {
	files, err := listMarkdownFiles()
	if err != nil {
		log.Fatalf("Error listing markdown files: %v", err)
	}

	if len(files) == 0 {
		log.Fatal("No markdown files found in the current directory.")
	}

	selectedFile, err := selectMarkdownFile(files)
	if err != nil {
		log.Fatal(err)
	}

	outputFilename, err := convertToHtmlWithFeatures(selectedFile)
	if err != nil {
		log.Fatalf("Error converting and saving file: %v", err)
	}

	openFile, err := confirmOpenFile()
	if err != nil {
		log.Fatal(err)
	}

	if openFile {
		err = openFileInBrowser(outputFilename)
		if err != nil {
			// Here we just log the error and continue, as it's not critical
			log.Printf("Error opening file in browser: %v", err)
		}
	}

	fmt.Println("Process completed successfully.")
}
