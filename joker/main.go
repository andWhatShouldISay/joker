package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	get := func(url string) ([]byte, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return ioutil.ReadAll(resp.Body)
	}

	unmarshalJoke := func(bytes []byte) (string, error) {
		type respJsonBody struct {
			Joke string `json:"value"`
		}
		var joke respJsonBody
		err := json.Unmarshal(bytes, &joke)
		if err != nil {
			return "", err
		}
		return joke.Joke, nil
	}

	flag.Parse()
	mode := flag.Arg(0)
	if mode == "random" {
		randomJoke := func() {
			bytes, err := get("https://api.chucknorris.io/jokes/random")
			if err != nil {
				log.Fatal(err)
			}
			joke, err := unmarshalJoke(bytes)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(joke)
		}
		randomJoke()
	} else if mode == "dump" {
		dump := func() {
			n := 5
			if flag.Arg(1) == "-n" {
				nval, err := strconv.ParseInt(flag.Arg(2), 10, 64)
				if err != nil {
					usage()
					return
				} else {
					n = int(nval)
				}
			}

			bytes, err := get("https://api.chucknorris.io/jokes/categories")
			if err != nil {
				log.Fatal(err)
			}
			var categories []string

			err = json.Unmarshal(bytes, &categories)
			if err != nil {
				log.Fatal(err)
			}

			getJokesByCategory := func(cat string, n int) []string {
				jokes := make([]string, n)
				for i := 0; i < n; i++ {
					bytes, err := get("https://api.chucknorris.io/jokes/random")
					if err != nil {
						log.Fatal(err)
					}
					jokes[i], err = unmarshalJoke(bytes)
					if err != nil {
						log.Fatal(err)
					}
				}
				return jokes
			}

			for _, cat := range categories {
				jokes := strings.Join(getJokesByCategory(cat, n), "\n")
				file, err := os.Create(cat + ".txt")
				if err != nil {
					log.Fatal(err)
				}
				_, err = file.Write([]byte(jokes))
				if err != nil {
					log.Fatal(err)
				}
			}

		}
		dump()
	} else {
		usage()
	}
}

func usage() {
	fmt.Println(`Usage:
joker random
joker dump
joker dump -n <number of jokes>`)
}
