package main

import (
	"github.com/andrebq/gas"
	"io/ioutil"
	"strings"
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type TrieNode struct {
	terminator bool
	children [26]*TrieNode
}

type Word struct {
	start int
	length int
	step int
}

var dictionary *TrieNode
var content []byte
var foundWords map[string][]*Word

func main() {
	loadDictionary();
	
	if len( os.Args ) < 2 {
		fmt.Println("Usage: code-finder [file to search]")
		os.Exit(0)
	}

	fmt.Println("Processing...")
	
	filename := os.Args[1]
	unfilteredContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	
	content = make([]byte, 0, len(unfilteredContent))
	for _, b := range unfilteredContent {
		if b >= 'A' && b <= 'Z' {
			content = append(content, b)
		} else if b >= 'a' && b <= 'z' {
			content = append(content, b - 'a' + 'A')
		}
	}
	
	foundWords = make(map[string][]*Word)
	
	//	Gather words.
	for step := -1000; step < 1000; step++ {
		if step < 0 || step > 1 {
			examineStep(step)
		}
	}
	
	writeTemplate()
}

func writeTemplate() {
	jsonWords := ""
	nWords := len(foundWords)
	wid := 0
	for k, v := range foundWords {
		if wid % 100 == 0 {
			fmt.Printf( "compiling: %d / %d\n", wid, nWords );
		}
		wid++
		
		jsonWords += fmt.Sprintf("\"%s\": [\n", k);
		for _, w := range v {
			jsonWords += fmt.Sprintf("\t[ %d, %d, %d ],\n", w.start, w.step, w.length)
		}
		jsonWords += "],\n"
	}
	
	template, err := gas.ReadFile("github.com/thingalon/code-finder/template.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	
	template = bytes.Replace(template, []byte("[filename]"), []byte(os.Args[1]), 1)
	template = bytes.Replace(template, []byte("[content]"), content, 1)
	template = bytes.Replace(template, []byte("[words]"), []byte(jsonWords), 1)
	
	outfile := os.Args[1] + ".html"
	
	ioutil.WriteFile(outfile, template, 0644)
	
	fmt.Println("Written to " + outfile)
}

func examineStep(step int) {
	started := 0
	done := make(chan bool, 100)

	for i := 0; i < len(content); i++ {
		go examineStepPos(done, step, i)
		started++
	}
	
	for started > 0 {
		<-done
		started--
	}
}

func examineStepPos(done chan bool, step int, start int) {
	length := 0
	cursor := dictionary
	pos := start

	for pos >= 0 && pos < len(content) {
		b := content[pos]
		index := int(b) - int('A')
		if cursor.children[index] == nil {
			done <- true
			return
		}
		
		cursor = cursor.children[index]
		length += 1
		
		if length >= 3 && cursor.terminator {
			word := getWord(step, start, length)
			if _, ok := foundWords[word]; ! ok {
				foundWords[word] = make([]*Word, 0, 1)
			}
			
			if len( foundWords[word] ) < 50 {
				foundWords[word] = append(foundWords[word], &Word{
					step: step,	
					start: start,
					length: length,
				})
			}
		}
		
		pos += step
	}
	
	done <- true
}

func getWord(step int, pos int, length int) string {	
	word := make([]byte, length, length)

	for i := 0; i < length; i++ {
		word[i] = content[pos]
		pos += step
	}

	return string(word)
}

func loadDictionary() {
	file, err := gas.Open("github.com/thingalon/code-finder/dictionary.txt")
	if err != nil {
		fmt.Println("Couldn't open dictionary")
		os.Exit(0)
	}
	defer file.Close()
	
	dictionary = &TrieNode{};	
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(strings.ToUpper(scanner.Text()))
		addToDictionary(line);
	}
}

func addToDictionary(line []byte) {
	cursor := dictionary
	for _, b := range line {
		index := int(b) - int('A');
		if cursor.children[index] == nil {
			cursor.children[index] = &TrieNode{}
		}
		
		cursor = cursor.children[index]
	}
	
	cursor.terminator = true
}

func wordInDictionary(word []byte) bool {
	cursor := dictionary
	for _, b := range word {
		index := int(b) - int('A')
		if cursor.children[index] == nil {
			return false
		}
		
		cursor = cursor.children[index]
	}
	
	return cursor.terminator
}

