package bitly

import (
	"github.com/gorilla/mux"
	"github.com/kurrik/twittergo"
	"github.com/kurrik/oauth1a"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"os"
	"fmt"
	"path/filepath"
)

func NewBitlyRouter(port string) *BitlyApi {
	var _port = ":" //let user declare port as "<int>"
	_port += port
	r := mux.NewRouter()
	fmt.Printf("\nRunning on port:%v", port)
	return &BitlyApi{Port: _port, Router: r}
}

type BitlyApi struct {
	Port string
	Twitter TwitterClient
  	Router *mux.Router
  	Route *mux.Route
}

type TwitterClient struct {
	Client *twittergo.Client
}
/*
	Handle your API path
 */
func(b *BitlyApi) HandleLinkHistory(path string) {
	route := GetBitlyRoute("v3/user/link_history")
	b.Router.HandleFunc(path, route.GetHistory).Methods("GET")
}

func(b *BitlyApi) HandleLinkSave(path string) {
	route := GetBitlyRoute("v3/user/link_save")
	b.Router.Handle(path, http.HandlerFunc(route.SaveLink)).Queries("longUrl","{longUrl}").Methods("GET")
}

func(b *BitlyApi) HandleUserClicks(path string) {
	route := GetBitlyRoute("v3/user/clicks")
	b.Router.HandleFunc(path, route.UserClicks).Methods("GET")
}

func(b *BitlyApi) HandleTweetBitlyUpdate(path string){
	route := GetTweetRoute("1.1/statuses/update.json")
	route.BitlyRoute = *GetBitlyRoute("v3/user/link_save")
	route.Client = b.Twitter.Client //Client is needed to make request.

	b.Router.HandleFunc(path, route.TweetAndBitly).Methods("POST")
}



func(b *BitlyApi) Listen() {
	if err := http.ListenAndServe(b.Port, b.Router); err != nil {
	log.Fatal(err)
	}
}

/*
	Set Twitter and Bitly Routes.
 */

func GetBitlyRoute(method string) (b *BitlyRoute){
	return &BitlyRoute{
		Scheme: "https",
		Host: "api-ssl.bitly.com",
		Method: method,
		Token: GetBitlyToken(),
		Param: "" } //Param is empty for queries in the route.
}

func GetTweetRoute(method string) (t *TwitterRoute){
	return &TwitterRoute{
		Scheme: "https",
		Host: "api.twitter.com",
		Method: method,
	}
}

/*
	Use OS environment variables to connect to twitter
 */

func(b *BitlyApi) GetTwitterCredEnv(){
	config := &oauth1a.ClientConfig{
		ConsumerKey: os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	twitterUser := oauth1a.NewAuthorizedConfig(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	b.Twitter.Client = twittergo.NewClient(config, twitterUser)
	return
}

/*
	Use a .txt file to connect to twitter
 */

func(b *BitlyApi) GetTwitterCredFile(filename string){
	absPath, _ := filepath.Abs(filename)
	creds, err := ioutil.ReadFile(absPath)
	if err != nil {
		return
	}
	lines := strings.Split(string(creds), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey: lines[0],
		ConsumerSecret: lines[1],
	}
	twitterUser := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	b.Twitter.Client = twittergo.NewClient(config, twitterUser)
	return
}

/*
	Verify you are connected to twitter.
 */

func(b *BitlyApi) VerifyTwitterCred(){

	var (
		err error
		client *twittergo.Client
		req *http.Request
		resp *twittergo.APIResponse
		user *twittergo.User
	)
	client = b.Twitter.Client
	req, err = http.NewRequest("GET", "/1.1/account/verify_credentials.json", nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	user = &twittergo.User{}
	err = resp.Parse(user)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\n---- Twitter Creds ----\n")
	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
}
