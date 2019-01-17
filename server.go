package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type questionsJSON struct {
	Questions []question `json:"questions,omitempty"`
}

type question struct {
	ID            int    `json:"id,omitempty"`
	IMG           string `json:"img,omitempty"`
	CorrectAnswer string `json:"correctAnswer,omitempty"`
	Hint          string `json:"hint,omitempty"`
}

var username = getEnv("GH_USER", "")
var passwd = getEnv("GH_TOKEN", "")
var secretsURL = getEnv("SECRETS_URL", "")

var refreshCRON = getEnv("REFRESH_CRON", "0 30 * * * *")
var port = getEnv("PORT", ":9287")

var cachedQuestions []question

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getQuestions() []question {

	client := &http.Client{}
	req, err := http.NewRequest("GET", secretsURL, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	result := questionsJSON{}
	err = json.Unmarshal(bodyText, &result)

	if err != nil {
		log.Fatal(err)
	}

	return result.Questions
}

func getCachedQuestions() []question {
	if cachedQuestions == nil {
		cachedQuestions = getQuestions()
	}

	return cachedQuestions
}

func main() {
	c := cron.New()
	c.AddFunc(refreshCRON, func() {
		fmt.Println("Update Cache")
		cachedQuestions = getQuestions()
	})
	c.Start()

	r := gin.Default()
	r.GET("/questions", func(c *gin.Context) {
		response := []question{}

		for index, q := range getCachedQuestions() {
			response = append(response, question{
				ID:   index,
				IMG:  q.IMG,
				Hint: q.Hint,
			})
		}

		c.JSON(200, response)
	})

	r.GET("/questions/:id/check/:answer", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		answer := c.Param("answer")

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		question := getCachedQuestions()[id-1]

		if strings.ToLower(question.CorrectAnswer) == strings.ToLower(answer) {
			c.JSON(200, gin.H{"correct": true})
		} else {
			c.JSON(200, gin.H{"correct": false})
		}
	})

	r.Run(port)
}
