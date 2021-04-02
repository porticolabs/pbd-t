package main 

import (
	"fmt"
	"encoding/json"
	"math/rand"
	
	log "github.com/sirupsen/logrus"

	"github.com/adjust/rmq/v3"
	"github.com/dghubble/go-twitter/twitter"
)

type TweetConsumer struct {
	name   string
	count  int
}

func NewConsumer(tag int) *TweetConsumer {
	return &TweetConsumer{
		name:   fmt.Sprintf("%s-consumer-%d", redisWorkerName, tag),
	}
}

func (consumer *TweetConsumer) Consume(delivery rmq.Delivery) {
    var tweet *twitter.Tweet
	
	log.WithFields(log.Fields{
		"tweetConsumer": consumer.name,
	}).Debug("Got new tweet from Redis Queue")

    if err := json.Unmarshal([]byte(delivery.Payload()), &tweet); err != nil {
        log.WithFields(log.Fields{
			"tweetConsumer": consumer.name,
		}).Warn("Couldn't unmarshal the json")
    
		if err := delivery.Reject(); err != nil {
            log.WithFields(log.Fields{
				"tweetConsumer": consumer.name,
			}).Error("Couldn't reject the queue tweet")
        }
    
		return
    }

    // perform task

	// Get the tweet fullText
	tweetText := getTweetText(tweet)
    log.WithFields(log.Fields{
		"tweetConsumer": consumer.name,
		"tweetText":    tweetText,
		"tweetID":    tweet.ID,
		"tweetUser":    tweet.User.ScreenName,
		}).Debug("Processing Tweet")
		
	// Analizes the tweet text to identify the request types (movie, book, serie, etc.)
	requestTypes := analizeTweetText(tweetText)
	log.WithFields(log.Fields{
		"tweetConsumer": consumer.name,
		"tweetText":    tweetText,
		"tweetID":    tweet.ID,
		"tweetUser":    tweet.User.ScreenName,
		"tweetRequestTypes": requestTypes,
		}).Debugf("Found Request Types: %s",requestTypes)
	
	// Get a recommendation for said request tipes
	recommendationRequest := &RecommRequest{}
	recommendationRequest.Tweet.ID      = tweet.IDStr
	recommendationRequest.Tweet.Text    = tweetText
	recommendationRequest.Tweet.User    = tweet.User.ScreenName
	recommendationRequest.Tweet.Request = requestTypes

	recommendations, recommendationsError := generateRecommendations(*recommendationRequest)
	
	if recommendationsError != nil {
		log.WithFields(log.Fields{
			"tweetConsumer": consumer.name,
			"tweetText":    tweetText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
			"tweetRequestTypes": requestTypes,
			"error": recommendationsError,
			}).Error("Couldn't get a recommendation")
		if err := delivery.Reject(); err != nil {
			log.WithFields(log.Fields{
				"tweetConsumer": consumer.name,
			}).Error("Couldn't reject the queue message")
		}
		log.WithFields(log.Fields{
			"tweetConsumer": consumer.name,
			"tweetText":    tweetText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
			"tweetRequestTypes": requestTypes,
			}).Debug("Queue message rejected")
		return
	}
	
	log.WithFields(log.Fields{
		"tweetConsumer": consumer.name,
		"tweetText":    tweetText,
		"tweetID":    tweet.ID,
		"tweetUser":    tweet.User.ScreenName,
		"tweetRequestTypes": requestTypes,
		"recommendation": recommendations[len(recommendations)-1].Recommendation.Text,
		}).Debug("Got a recommendation")
		
	// Reply on Twitter
	if !twitterSample {
		for _, recommendation := range recommendations {
			_, tweetError := replyTweet(tweet, 
				recommendation.Recommendation.Text, 
				recommendation.Recommendation.MediaURL)
			if tweetError != nil {
				log.WithFields(log.Fields{
					"tweetConsumer": consumer.name,
					"tweetText":    tweetText,
					"tweetID":    tweet.ID,
					"tweetUser":    tweet.User.ScreenName,
				}).Error("Couldn't reply tweet")
				return
			} 
		}
		
		_, retweetError := retweetTweet(tweet)
		if retweetError != nil {
			log.WithFields(log.Fields{
				"tweetConsumer": consumer.name,
				"tweetText":    tweetText,
				"tweetID":    tweet.ID,
				"tweetUser":    tweet.User.ScreenName,
			}).Error("Couldn't retweet tweet")
			return
		} 
		log.WithFields(log.Fields{
			"tweetConsumer": consumer.name,
			"tweetText":    tweetText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
			"tweetRequestTypes": requestTypes,
			"recommendation": recommendations[len(recommendations)-1].Recommendation.Text,
			}).Info("Tweet Replied and Retweeted")
	} 
	
    if err := delivery.Ack(); err != nil {
        log.WithFields(log.Fields{
			"tweetConsumer": consumer.name,
			"tweetText":    tweetText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
		}).Error("Couldn't acknowledge queue message")
    }
	log.WithFields(log.Fields{
		"tweetConsumer": consumer.name,
		"tweetText":    tweetText,
		"tweetID":    tweet.ID,
		"tweetUser":    tweet.User.ScreenName,
	}).Debug("Queue message acked")
}

// Tweet that are more than 144 characters long, use a secundary
//   Tweet structure called 'ExtendedTweet' which contains the
//   complete text, and not just the trucated one.
// This functions get returns the tweet complete text wheter if
//   it's extended or not.
func getTweetText (tweet *twitter.Tweet) (string) {
	if (tweet.Truncated) {
		log.WithFields(log.Fields{
			"tweetText":    tweet.ExtendedTweet.FullText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
		}).Debug("Text truncated, checking Extended Tweet")
		return tweet.ExtendedTweet.FullText
	} else {
		return tweet.Text
	}
}

func generateRecommendations(request RecommRequest) ([]RecommResponse, error){
	var recommendations []RecommResponse

	if len(request.Tweet.Request) == 0 {
		request.Tweet.Request = []string{keywordsTypeList[rand.Intn(len(keywordsTypeList))]}
	}
	
	recommendationResponse, recommError := getRecommendation(request)
	
	if recommError != nil {
		return recommendations, recommError
	}
	
	recommendations = append(recommendations, recommendationResponse)

	if !recommendationResponse.Fulfilled {
		// The recommendation API doesn't has anything to recommend
		//  for the given request type.
		log.WithFields(log.Fields{
			"tweetText":    request.Tweet.Text,
			"tweetID":    request.Tweet.ID,
			"tweetUser":    request.Tweet.User,
			"tweetRequestTypes": request.Tweet.Request, 
			}).Debug("Couldn't fulfill recommendation request, getting a random recommendation")
		
		request.Tweet.Request = []string{keywordsTypeList[rand.Intn(len(keywordsTypeList))]}
		
		randRecommendationResponse, randRecommError := getRecommendation(request)
		
		if randRecommError != nil {
			return recommendations, randRecommError
		}
		recommendations = append(recommendations, randRecommendationResponse)
	}
	return recommendations, nil
}