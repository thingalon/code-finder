Code Finder
===========

This hacky little goscript is for finding hidden messages in any text, made from equidistant letters, ala "the bible code".

If your go environment is properly set up, you can get it with:

	go get github.com/thingalon/code-finder

Then to run it: 

	$GOPATH/bin/code-finder <file to process>

If you would like to add custom words to the dictionary (by default, it just uses a standard dict file):
	
	vim $GOPATH/src/github.com/thingalon/code-finder/dictionary.txt
	
It will output an html doc that you can then open to browse the results. You can set your crossword width, offset in the source document, and click on words to hilight them.

Known bug: If you select a word that wraps around the edges of the crossword, the line drawn for the word will be wrong. Only words that fit in your current crossword width will look correct.
