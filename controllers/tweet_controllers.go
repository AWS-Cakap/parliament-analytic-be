package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"parliament-analytic-be/config"
	"parliament-analytic-be/models"

	"github.com/gin-gonic/gin"
)

type TweetAPIResponse struct {
	Data []struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
		AuthorID  string `json:"author_id"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
}

func AnalyzeTweetSentiment(text string) string {
	payload := map[string]string{"text": text}
	jsonValue, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:8000/analyze", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Error calling sentiment API:", err)
		return "unknown"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	return result["sentiment"]
}

func CrawlTweets(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		query = "Gerindra"
	}

	url := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%s lang:id -is:retweet&tweet.fields=created_at,author_id&expansions=author_id&user.fields=username&max_results=15", query)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to call Twitter API"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response TweetAPIResponse
	json.Unmarshal(body, &response)

	// Map author_id to username
	userMap := make(map[string]string)
	for _, user := range response.Includes.Users {
		userMap[user.ID] = user.Username
	}

	for _, tweet := range response.Data {
		tweetTime, _ := time.Parse(time.RFC3339, tweet.CreatedAt)
		username := userMap[tweet.AuthorID]
		sentiment := AnalyzeTweetSentiment(tweet.Text)

		t := models.Tweet{
			Username:  username,
			Text:      tweet.Text,
			CreatedAt: tweetTime,
			Sentiment: sentiment,
		}

		config.DB.Create(&t)
	}

	c.JSON(200, gin.H{"message": "Tweets saved"})
}
