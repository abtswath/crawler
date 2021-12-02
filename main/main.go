package main

import (
	"fmt"
	"regexp"
)

func main() {
	regex := regexp.MustCompile("[a-zA-Z]+")
	result := regex.FindAllString("asdfouq8924uerquiwer0q9ri349i43jiwrejkvd323", -1)
	for _, s := range result {
		fmt.Println(s)
		fmt.Println(s[1:len(s)-1])
	}
}
