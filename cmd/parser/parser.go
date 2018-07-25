package main

import (
	"encoding/json"
	"log"

	"github.com/mjudeikis/sample-ini-parser/pkg/parser"
)

func main() {

	sampleInput := `
		key=test1
		[Data]
		data=datacontent
		data2=datacontent2
		[Data2]
		data=datacontent
	`

	parsedINIFile := parser.Parse("sample.ini", sampleInput)
	prettyJSON, err := json.MarshalIndent(parsedINIFile, "", "   ")

	if err != nil {
		log.Println("Error marshalling JSON:", err.Error())
		return
	}

	log.Println(string(prettyJSON))

}
