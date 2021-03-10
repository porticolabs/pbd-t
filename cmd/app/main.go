package main

import (
	"os"
    "os/signal"
    "syscall"
    "strconv"
	
    log "github.com/sirupsen/logrus"
)

var (
    version             string // version number
    sha1ver             string // sha1 revision used to build the program
    buildTime           string // when the executable was built
    twitterCreds        TwitterCredentials // Twitter login data
    twitterHashtag      string // Twitter Hashtag to follow
    twitterSample       bool // Use a sample stream for testing purposes
    redisCreds          RedisCredentials // Redis connection data
    redisQueue          string // Redis Queue to publish curated tweets
    redisWorkerName     string // Redis consumer name prefix
    redisConsumers      int // Number of tweet consumers to create 
    recommProtocol      string // 
    recommHost          string // 
    recommPort          string // 
)

func init() {
    // Get the settings from the environment variables
    if os.Getenv("ENVIRONMENT") == "prod" {
        log.SetFormatter(&log.JSONFormatter{})
      } 

    // Getting log level settings
    switch os.Getenv("LOG_LEVEL") {
        case "DEBUG":
            log.SetLevel(log.DebugLevel)
            log.Warn("Log level set to DEBUG")
        case "WARN":
            log.SetLevel(log.WarnLevel)
            log.Warn("Log level set to WARN")
        case "ERROR":
            log.SetLevel(log.ErrorLevel)
        case "FATAL":
            log.SetLevel(log.FatalLevel)
    }

    // Getting the Twitter credentials from env
	twitterCreds = TwitterCredentials{
        AccessToken:       os.Getenv("ACCESS_TOKEN"),
        AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
        ConsumerKey:       os.Getenv("CONSUMER_KEY"),
        ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
    }

    // Getting twitter stream settings
    twitterHashtag = os.Getenv("TWITTER_HASHTAG")
    
    // Run twitter interaction avoiding real answers
    twitterSample, _ = strconv.ParseBool(os.Getenv("TWITTER_SAMPLE"))
    
    // Getting the Redis credentials from env
	redisCreds = RedisCredentials{
        Host:     os.Getenv("REDIS_HOST"),
        Protocol: os.Getenv("REDIS_PROTOCOL"),
    }
    redisQueue = os.Getenv("REDIS_QUEUE")
    redisConsumers, _ = strconv.Atoi(os.Getenv("REDIS_TWEET_CONSUMERS"))
    redisWorkerName = os.Getenv("REDIS_WORKER_NAME")

    recommProtocol = os.Getenv("RECOMM_PROTOCOL")
    recommHost = os.Getenv("RECOMM_HOST")
    recommPort = os.Getenv("RECOMM_PORT")
  }

func main() {
    // Startup Information
    log.Infof("Initiating PBD (Twitter) v%s", version)
    log.Infof(" * Commit Hash: %s", sha1ver)
    log.Infof(" * Build Date: %s", buildTime)
	
	//log.Debug("Signing in to Twitter.")
	//loginToTwitter(&twitterCreds)
    
	log.Debug("Connecting to Redis")
    loginToRedis(redisCreds) // Login into Redis service

    log.Debug("Opening Redis Queue")
    openQueue(redisQueue) // Open the Redis queue to send the tweets (It'll created if it doesn't exist)

	err := startConsumingQueue()
	if err != nil {
		log.Warn("Error starting Redis Queue consuming")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals // wait for signal
	go func() {
		<-signals // hard exit on second signal (in case shutdown gets stuck)
		os.Exit(1)
	}()

	<-stopConsumingQueue() // wait for all Consume() calls to finish
}