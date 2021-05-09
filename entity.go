package main

import "golang.org/x/text/language"

type RequestModel struct {
	from       language.Tag
	to         language.Tag
	fromPhrase string
}
