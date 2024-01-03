package main

import (
	"encoding/json"
	"log"
	"os"
	wc "wc-tool/internal"
)

func main() {
	result := wc.Process(os.Args)

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(jsonResult)

}
