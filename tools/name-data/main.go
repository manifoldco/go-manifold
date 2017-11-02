package main

import (
	"bufio"
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func execTemplate(tmpl *template.Template, data interface{}) *bytes.Buffer {
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, data)
	if err != nil {
		log.Fatal(err)
	}

	return buf
}

func writeTemplate(fileName string, tmpl *template.Template, data interface{}) {
	buf := execTemplate(tmpl, data)

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		log.Fatal("Error creating dir:", err)
	}

	fd, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	fd.Write(formatted)
}

func readlines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	infile := os.Args[1]
	outfile := os.Args[2]

	fname := filepath.Base(infile)
	fname = fname[:len(fname)-4]

	lines, err := readlines(infile)
	if err != nil {
		log.Fatal("Error reading input file:", err)
	}

	sliceTmpl := template.Must(template.ParseFiles("tools/name-data/slice.tmpl"))
	writeTemplate(outfile, sliceTmpl, struct {
		Name  string
		Lines []string
	}{
		Name:  strings.Title(fname),
		Lines: lines,
	})
}
