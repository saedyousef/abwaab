package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/saedyousef/abwaab-task/auth"
	"github.com/saedyousef/abwaab-task/models"
	"github.com/saedyousef/abwaab-task/tweeters"
)

type CreateTweetInput struct {
	Body  string `json:"body" binding:"required"`
	UserID uint64 `json:"user_id"`
}

// This api call a function that uses Twitter Client to get tweets.
func SearchTweets(c *gin.Context) {

	// Validate input
	url := c.Request.URL.Query()
	if url["query"] == nil || len(url["query"][0]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL param 'query' is required"})
	  	return
	}

	tweets := tweeters.SearchTweets(url["query"][0])

	// Extracting the token to get the userid.
	tokenAuth, err := auth.ExtractTokenMetadata(c.Request)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	for _, tweet := range tweets {
		go SaveTweet(tweet, tokenAuth.UserId)
	}

	c.JSON(http.StatusOK, tweets)
}


func SaveTweet(body string, userId uint64) {
	tweet := models.Tweet{Body: body, UserID: userId}
	models.DB.Create(&tweet)
}

func CreateTweet(c *gin.Context) {
	var input CreateTweetInput

		// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	// Extracting the token to get the userid.
	tokenAuth, err := auth.ExtractTokenMetadata(c.Request)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	
	tweet := models.Tweet{Body: input.Body, UserID: tokenAuth.UserId}
	models.DB.Create(&tweet)

	c.JSON(http.StatusCreated, gin.H{"data": tweet})
}

// Getting tweets for logged in user.
func GetUserTweets(c *gin.Context) {
	var tweets []models.Tweet

	err := models.DB.Where("user_id = ?", 142).Order("id desc").Find(&tweets).Error
	if err != nil {
		c.JSON(http.StatusBadRequest,  gin.H{"error": err})
		return
	}

	if len(tweets) == 0 {
		c.JSON(http.StatusNotFound,  gin.H{"error": "No tweets found for this user."})
		return
	}


	c.JSON(http.StatusOK,  gin.H{"data": tweets})
}
