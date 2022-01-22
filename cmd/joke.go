/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResult struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

// jokeCmd represents the joke command
var jokeCmd = &cobra.Command{
	Use:   "joke",
	Short: "输出一条笑话",
	Long:  `输出一条笑话`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getRandomJokeWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}
	},
}

func init() {
	rootCmd.AddCommand(jokeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	jokeCmd.PersistentFlags().String("term", "", "要查询的关键字")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}

	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	fmt.Println(string(joke.Joke))
}

func getRandomJokeWithTerm(jokeTerm string) {
	total, results := getJokeDataWithTerm(jokeTerm)
	randomiseJokeList(total, results)
}

func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(
		http.MethodGet, //method
		baseAPI,        //url
		nil,            //body
	)

	if err != nil {
		log.Printf("Could not request a dadjoke. %v", err)
	}

	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}

	return responseBytes
}

func getJokeDataWithTerm(jokeTerm string) (totalJokes int, jokeList []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)

	jokeListRaw := SearchResult{}

	if err := json.Unmarshal(responseBytes, &jokeListRaw); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	jokes := []Joke{}
	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	return jokeListRaw.TotalJokes, jokes
}

func randomiseJokeList(length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())

	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("No jokes found with this term")
		fmt.Println(err.Error())
	} else {
		randomNum := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNum].Joke)
	}
}
