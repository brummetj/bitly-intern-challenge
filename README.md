#### Bitly API internship challenge 2018

#### Joshua Brummet

## Introduction 

For your internship challenge, I developed a library that acts as both a Router and Wrapper to both the Bitly API and Twitter API, written in go. With this library, you can submit a new tweet with some text. If there is a URL present in your string, it will replace this URL with a bitly link and submit the tweet to your twitter account. 

This library was built to make quick responses to twitter with bitly links.

For example... Let us say someone uses Twitter as a marketing tool, and they have items they sell with URI's that you would describe as a 'Deep link'. This library API could make it so someone can quickly market their item as a Bitly link on their twitter account. This person can utilize all the valuable information you guys provide fast and easy with one HTTP request to their own endpoint!! :)

If there is no URL present in the form value then it will simply just submit a tweet to your account. 

Built upon this functionality, this library also includes the following endpoints.

* `/v3/user/link_save`
* `/v3/user/link_history`
* `/v3/user/clicks`

The twitter endpoints used are 

* `/1.1/statuses/update.json`
* `/1.1/account/verify_credentials.json`

In order to use the Twitter functionality with the library, this will require you to set up a quick app and get your user credentials. This will take only take a moment to set up. Just go to this [link](https://apps.twitter.com/)

Once you have created your app you will need to get the following tokens

* `TWITTER_ACCESS_TOKEN`
* `TWITTER_ACCESS_TOKEN_SECRET` 
* `TWITTER_CONSUMER_KEY`
* `TWITTER_CONSUMER_SECRET`

---
## How to Run 

To run the library you can either use the folder I have provided with you guys, or simply just use the google go package I created. This package will include all the dependencies so there is no need to install the library dependencies when testing. 

* `go get -u github.com/brummetj/bitly-intern-challenge`

If you decide to use the folder as a package here is a list of dependencies you will need to install.

* `go get -u github.com/gorilla/mux`
* `go get -u github.com/kurrik/twittergo`
* `go get -u github.com/kurrik/oauth1a`
* `go get -u mvdan.cc/xurls`

`mkdir` A new directory, cd inside and let's call our new bitly app.. `app.go`

### Authentication 

You will want to `touch` a new text file that will hold your twitter credentials. I call mine `CREDENTIALS.txt`.
NOTE: I have provided a function so you can just export the tokens as OS variables if that is easier for you. 

Inside the `CREDENTIALS.txt` insert the tokens in the following order. 
```
<TWITTER_CONSUMER_KEY>
<TWITTER_CONSUMER_SECRET
<TWITTER_ACCESS_TOKEN>
<TWITTER_ACCESS_TOKEN_SECRET>
```
If you decide to use OS variables just run for each variable in the command line such as.
`export TWITTER_CONSUMER_KEY=<YOUR_TOKEN>`

Your folder structure should look like this.

```
-dir
   -- app.go
   -- CREDENTIALS.txt
```

For the bitly API I use an exported variable for my credentials to access your API. 

Simply in the command line do
`export BITLY_ACCESS_TOKEN=<YOUR_TOKEN>`

This will allow each bitly endpoint to be successful.

## Using

Inside your app.go you can just copy this following code.

```
package main

import "github.com/brummetj/bitly-intern-challenge"

func main(){
    
    r := bitly.NewBitlyRouter("5000")

    r.GetTwitterCredFile("./CREDENTIALS.txt")
     
  //r.GetTwitterCredEnv()
  
    r.VerifyTwitterCred()
    
    r.HandleLinkHistory("/link_history")
    
    r.HandleLinkSave("/link_save")
    
    r.HandleUserClicks("/user_clicks")
    
    r.HandleTweetBitlyUpdate("/tweet")
    
    r.Listen()
}
```
To Run please type `go run app.go`

You will see a message on which port you running. 

Open up Postman and you can start making requests to your localhost endpoints which will respond to the Bitly and Twitter API.

NOTE: if you plan to you use the source code as your package dependency, make sure you install the listed dependencies. You will also want to include the source code in a dir outside of app.go
```
package main

import "../bitly-intern-challenge"

```
And the folder structure will look like 

```
-dir
  -src
    -- app.go
    -- CREDENTIALS.txt
  -bitly-intern-challenge
```

#### Bitly Router & Listen

The bitly router is built on top of gorilla's Mux package. Declaring `bitly.NewBitlyRouter("5000")` Will instantiate a new router that can be used to declare each handler and its endpoint. This function will take an argument as a port number. This is so you can have multiple local servers running on different ports. 

The BitlyRouter has to always be the first variable declared in order to use each handler

`Listen()` is needed at the end of each handler call. This will make the API live and ready to use!! 

#### Twitter Creds

`r.GetTwitterCredFile("./CREDENTIALS.txt")` will simply look for the `CREDENTIALS.txt` file in your dir, and authenticate you to Twitter to make calls to their API. This function is needed before making and handler endpoints to Twitter. 

alternatively if you can use `r.GetTwitterCredEnv()` if you declare OS variables. 

`r.VerifyTwitterCred()` Will verify if you have successfully connected to Twitter or not. 

If you don't have a `CREDENTIALS` file in the dir and call `GetTwitterCredFile` you will get a compile error.

## Handlers

Each handler will take an endpoint, and allow you to make calls to the endpoint depending on the HTTP request. Most of the handlers declared in the library are just GET requests. The twitter / bitly handler is a POST request. They all will respond in JSON.

All Tests are used with Postman, feel free to `CURL` or use the web browser for getting requests.

* `GET Request`
* `r.HandleLinkHistory("/link_history")`
  * This will retrieve a response from `v3/user/link_history` in the bitly API.
  
Test Response from `HandleLinkHistory("/link_history")`

* `GET localhost:5000/link_history`
```
{
    "data": {
        "link_history": [
            {
                "has_link_deeplinks": false,
                "archived": false,
                "user_ts": 1522789253,
                "title": "Mr. E by Metroplane | Free Listening on SoundCloud",
                "created_at": 1522789253,
                "tags": [],
                "modified_at": 1522789253,
                "campaign_ids": [],
                "private": true,
                "aggregate_link": "http://bit.ly/20KUVJU",
                "long_url": "https://soundcloud.com/metroplane/metroplane-mr-e",
                "client_id": "a5e8cebb233c5d07e5c553e917dffb92fec5264d",
                "link": "http://bit.ly/2H6JLwn",
                "is_domain_deeplink": false,
                "encoding_user": {
                    "login": "o_5hfavnimd",
                    "display_name": "Joshua Brummet",
                    "full_name": "Joshua Brummet"
                }
            },
```

* `GET Request`
* `Query: longUrl`
* `r.HandleLinkSave("/link_save")`
  * This will retrieve a response from `v3/user/link_save` in the bitly API.
  
Test Response from `HandleLinkSave("/link_save")`

* `GET localhost:5000/link_save?longUrl=https://codeburst.io/6-interesting-apis-to-check-out-in-2018-5d6830063f29`
```
{
    "status_code": 200,
    "data": {
        "link_save": {
            "link": "http://bit.ly/2EkT3Sw",
            "aggregate_link": "http://bit.ly/2Cqg1eC",
            "long_url": "https://codeburst.io/6-interesting-apis-to-check-out-in-2018-5d6830063f29",
            "new_link": 1,
            "user_hash": "2EkT3Sw"
        }
    },
    "status_txt": "OK"
}
```
* `GET Request`
* `r.HandleUserClicks("/user_clicks")`
  * This will retrieve a response from `v3/user/clicks` in the bitly API.
  
Test Response from `HandleUserClicks("/user_clicks")`

* `GET localhost:5000/user_clicks`
```
{
    "status_code": 200,
    "data": {
        "days": 7,
        "day_start": 0,
        "clicks": [
            {
                "clicks": 0,
                "day_start": 1522728000
            },
            {
                "clicks": 0,
                "day_start": 1522641600
            },
            {
                "clicks": 0,
                "day_start": 1522555200
            },
            {
                "clicks": 0,
                "day_start": 1522468800
            },
            {
                "clicks": 3,
                "day_start": 1522382400
            },
            {
                "clicks": 1,
                "day_start": 1522296000
            },
            {
                "clicks": 0,
                "day_start": 1522209600
            }
        ]
    },
    "status_txt": "OK"
}
```
* `POST Request`
* `r.HandleTweetBitlyUpdate("/tweets")`
  * This will retrieve a response from `/1.1/statuses/update.json` in the Twitter API.
  
Test Response from `HandleTweetBitlyUpdate("/tweets")`

* `POST localhost:5000/tweets`
  * `x-www-form-urlencoded`
  * `Key = status`
  * `Value = Hey there ! https://golang.org/`
```
{
    "Id": 981257126097006592,
    "Tweet": "Hey there ! https://t.co/UUpoxF7Pcp",
    "User": "Joshua Brummet"
}
```
##### For simplicity purposes, I weeded out most of the twitter response to save the time of creating a model for twitters API response.

Here screenshot of my twitter account with the new bitly / twitter update! 

![alt text](https://github.com/brummetj/bitly-intern-challenge/blob/master/Screen%20Shot%202018-04-03%20at%206.05.28%20PM.png)

---
## Conclusion

With this API go wrapper you can quickly submit bitly links to twitter through the API router that I built with your own endpoints! 
