package main

import (
	"fmt"

	"github.com/maengsoo/learngo/dictionary/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	baseWord := "hello"

	dictionary.Add(baseWord, "First")

	err := dictionary.Update(baseWord, "Second")
	if err != nil {
		fmt.Println(err)
	}

	err2 := dictionary.Delete(baseWord)
	if err2 != nil {
		fmt.Println(err2)
	}

	_, err3 := dictionary.Search(baseWord)
	fmt.Println(err3)

	// word := "hello"
	// definition := "Greeting"
	// err := dictionary.Add(word, definition)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// hello, _ := dictionary.Search(word)
	// fmt.Println(hello)
	// err2 := dictionary.Add(word, definition)
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }
}
