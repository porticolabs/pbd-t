package main

import (
    "io/ioutil"
    "net/http"

    log "github.com/sirupsen/logrus"

    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
)

var (
    client *twitter.Client
    user *twitter.User
)

// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type TwitterCredentials struct {
    ConsumerKey       string
    ConsumerSecret    string
    AccessToken       string
    AccessTokenSecret string
}

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func getTwitterClient(creds *TwitterCredentials) (*twitter.Client, *twitter.User, error) {
    // Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
    config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
    // Pass in your Access Token and your Access Token Secret
    token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

    httpClient := config.Client(oauth1.NoContext, token)
    client := twitter.NewClient(httpClient)

    // Verify Credentials
    verifyParams := &twitter.AccountVerifyParams{
        SkipStatus:   twitter.Bool(true),
        IncludeEmail: twitter.Bool(true),
    }

    // we can retrieve the user and verify if the credentials
    // we have used successfully allow us to log in!
    user, _, err := client.Accounts.VerifyCredentials(verifyParams)
    if err != nil {
        return nil, nil, err
    }

    return client, user, nil
}

// Uses the twitter login data to authenticate with the 
//  Twitter API and stores the connection object in a
//  global variable
func loginToTwitter(creds *TwitterCredentials){
	var err error
	client, user, err = getTwitterClient(creds)
    if err != nil {
        log.Warn("Error getting Twitter Client")
        log.Error(err)
    }
	log.Info("Logged in as User: " + user.Name)
}

func retweetTweet(originTweet *twitter.Tweet) (*twitter.Tweet, error) {

    // Retweeting the original tweet
    retweetParams := &twitter.StatusRetweetParams{
        TrimUser: twitter.Bool(true),
    }
    retweet, retweetResponse, retweetError := client.Statuses.Retweet(originTweet.ID, retweetParams)
    
	log.WithFields(log.Fields{
        "tweetID":    originTweet.ID,
        "tweetUser":    originTweet.User.ScreenName,
        "tweetAPIResponse": retweetResponse.StatusCode,
    }).Debug("Tweet Retweet Status Code")
	
	return retweet, retweetError
}

func replyTweet(originTweet *twitter.Tweet, answerText string, mediaURL string) (*twitter.Tweet, error) {
	tweetParams := &twitter.StatusUpdateParams{
		InReplyToStatusID: originTweet.ID,
	}
        
	if mediaURL != "" {
		mediaID, mediaError := getMediaID(mediaURL)
		if mediaError != nil {
			return nil, mediaError
		}
		tweetParams.MediaIds = []int64{mediaID}
	}

	newTweet, resp, tweetError := client.Statuses.Update(answerText, tweetParams)
	log.WithFields(log.Fields{
        "tweetAnswer":    answerText,
        "tweetID":    originTweet.ID,
        "tweetUser":    originTweet.User.ScreenName,
        "tweetAPIResponse": resp.StatusCode,
    }).Debug("Tweet Reply Status Code")
	
	return newTweet, tweetError
}

func getMediaID (imageURL string) (int64, error) {
	image, imageError := getImage(imageURL)
	if imageError != nil {
        return 0, imageError
    }

	media, _, error := client.Media.Upload(image, "IMAGE")
    return media.MediaID, error
}

func getImage (imageURL string) ([]byte, error) {
    resp, getErr := http.Get(imageURL)
    if getErr != nil {
        return nil, getErr
    }
    defer resp.Body.Close()

    image, imgErr := ioutil.ReadAll(resp.Body)

    return image, imgErr
}