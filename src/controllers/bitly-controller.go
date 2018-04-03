package controllers

import (
 	 "../models"
 	 "../utility"
	"net/http"
  	"net/url"
	"fmt"
	"log"
	"encoding/json"
	"github.com/kurrik/twittergo"
	"mvdan.cc/xurls"
	"strings"
)

/*
	Controller Route Objects
 */

type BitlyRoute struct{
	Scheme string
	Method string
	Param string
	ParamTrue bool
	Host string
	Token string
}
type TwitterRoute struct {
	Param string
	Scheme string
	Host string
	Method string
	Body string
	BitlyRoute BitlyRoute
	Client *twittergo.Client
}
type Error struct {
	Error int
	StatusCode int16
	Message string
  	Code int64
}

/*
	Controller Routers to handle REST API request.
	Params: HTTP.responseWriter, HTTP.Request
	Returns: A JSON response
 */

func (b * BitlyRoute) GetHistory(w http.ResponseWriter, r *http.Request) {

	resp, err := GetBitlyResp(b)
	if err != nil {
		e := Error{StatusCode: 500, Message: "Oops! - please look at the GetHistory controller"}
		log.Fatal("Do: ", err)
		utility.RespondWithJson(w, http.StatusCreated, e)
		return
	}
	defer resp.Body.Close()
	var history models.HistoryModel
	if err := json.NewDecoder (resp.Body).Decode(&history); err != nil {
		log.Println(err)
	}
	utility.RespondWithJson(w, http.StatusCreated, history)
}

func (b * BitlyRoute) SaveLink(w http.ResponseWriter, r *http.Request){

	/*
		Requests for Saving a Link will need key = "longUrl"
	 */
	vals := r.URL.Query()
	longUrls, ok := vals["longUrl"]
	var u string
	if ok {
		if len(longUrls) >= 1 {
			u = longUrls[0]
		}
	}
	b.Param = fmt.Sprintf("longUrl=%s", u)
	resp, err := GetBitlyResp(b)
	if err != nil {
		e := Error{StatusCode: 500, Message: "Oops! - please look at the GetHistory controller"}
		log.Fatal("Do: ", err)
		utility.RespondWithJson(w, http.StatusCreated, e)
		return
	}
	defer resp.Body.Close()
	var savedLink models.LinkModel

	if err := json.NewDecoder (resp.Body).Decode(&savedLink); err != nil {
		log.Println(err)
	}

	utility.RespondWithJson(w, http.StatusCreated, savedLink)
}

func (b * BitlyRoute) UserClicks(w http.ResponseWriter, r *http.Request){

	resp, err := GetBitlyResp(b)
	if err != nil {
		e := Error{StatusCode: 500, Message: "Oops! - please look at the UserClicks controller"}
		log.Fatal("Do: ", err)
		utility.RespondWithJson(w, http.StatusCreated, e)
		return
	}
  	defer resp.Body.Close()
  	var clicks models.ClickModel

  	if err := json.NewDecoder (resp.Body).Decode(&clicks); err != nil {
    	log.Println(err)
  	}

  	utility.RespondWithJson(w, http.StatusCreated, clicks)

}
func (t * TwitterRoute) TweetAndBitly(w http.ResponseWriter, r *http.Request){

	/*
	 *	Parsing Form request for Text
	 */
	r.ParseForm()
	b := r.Form
	body := b["status"]
	k := strings.Join(body, " ")
	var bodyArray = strings.Fields(k)
	var bUrl = xurls.Relaxed().FindString(k)
	t.BitlyRoute.Param = fmt.Sprintf("longUrl=%s", bUrl)

	/*
	 *	Making a request to Bitly for the URL.
	 */
	bResp, err := GetBitlyResp(&t.BitlyRoute)
	if err != nil {
		e := Error{StatusCode: 500, Message: "Oops! - please look at the TweetAndBitly controller"}
		log.Fatal("Do: ", err)
		utility.RespondWithJson(w, http.StatusCreated, e)
		return
	}

	defer bResp.Body.Close()
	var savedLink models.LinkModel
	if err := json.NewDecoder (bResp.Body).Decode(&savedLink); err != nil {
		log.Println(err)
	}

	/*
	 *	Replace Url with bitly link
	 */
	var newBitlyLink = savedLink.Data.LinkSavedObject.AggregateLink
	for i := 0; i < len(bodyArray); i++ {
		if bodyArray[i] == bUrl{
			bodyArray[i] = newBitlyLink
		}
	}

	/*
	 * Make the request to post a tweet to your account.
	 */
	tweetBody := strings.Join(bodyArray, " ")
	t.Body = tweetBody
	tweetResp, err := GetTweetResp(t)
	if err != nil {
		e := Error{StatusCode: 500, Message: "Oops! - please look at the TweetAndBitly controller"}
		log.Fatal("Do: ", err)
		utility.RespondWithJson(w, http.StatusCreated, e)
	}
	tweet := &twittergo.Tweet{}
	err = tweetResp.Parse(tweet)
	if err != nil {
		if rle, ok := err.(twittergo.RateLimitError); ok {
			fmt.Printf("Rate limited, reset at %v\n", rle.Reset)
		} else if errs, ok := err.(twittergo.Errors); ok {
			for i, val := range errs.Errors() {
        		err := Error{
        		  	Error: i + 1,
				  	Code: val.Code(),
				  	Message: val.Message(),
					}
				fmt.Printf("Error #%v - ", i + 1)
				fmt.Printf("Code: %v ", val.Code())
				fmt.Printf("Msg: %v\n", val.Message())
				utility.RespondWithJson(w, http.StatusCreated, err)
				return
			}
		} else {
			fmt.Printf("Problem parsing response: %v\n", err)
			err := Error{
					Code: 500,
					Message: "Problem parsing response",
			}
			utility.RespondWithJson(w, http.StatusCreated, err)
			return
		}
		//os.Exit(1)
	}
	var tweetObject models.Tweet
	tweetObject.Id = tweet.Id()
	tweetObject.Tweet = tweet.Text()
	tweetObject.User = tweet.User().Name()
	utility.RespondWithJson(w, http.StatusCreated, tweetObject)
}

/*
	Get HTTP responses from the requested API.
 */

func GetBitlyResp(b *BitlyRoute) (*http.Response, error){
	route :=
		fmt.Sprintf("%s://%s/%s?access_token=%s&%s",
		b.Scheme,
		b.Host,
		b.Method,
		url.PathEscape(b.Token),
		url.PathEscape(b.Param))

	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil, err
	}
	return response, err
}

func GetTweetResp(t *TwitterRoute) (*twittergo.APIResponse, error){
	route :=
		fmt.Sprintf("%s://%s/%s",
			t.Scheme,
			t.Host,
			t.Method)
	data := url.Values{}
	data.Set("status", t.Body)
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", route, body)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := t.Client
	if client.User != nil {
		client.OAuth.Sign(req, client.User)
	} else {
		if err = client.Sign(req); err != nil {
			return nil, err
		}
	}
	var r *http.Response
	r, err = client.HttpClient.Do(req)
	resp := (*twittergo.APIResponse)(r)
	return resp, nil
}
