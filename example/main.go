package main

import (
	"fmt"
	"log"

	"github.com/blainsmith/infogram-go"
)

func main() {
	client := infogram.Client{
		Endpoint:  infogram.DefaultEndpoint,
		APIKey:    "VoyBH3SykNCgqcWD9CybuPxwVGFToUJ3",
		APISecret: "qUoyG18UrkC0XGbgf7vOfhVy7ddcXWTw",
	}

	themes, err := client.Themes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(themes)

	infographics, err := client.Infographics()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(infographics)

	infographic, err := client.Infographic(infographics[0].Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(infographic)
}
