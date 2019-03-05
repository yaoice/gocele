package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	version                  = "1.0"
	SERVICE_NAME             = "gocele"
	FLAG_KEY_SERVER_HOST     = "server.host"
	FLAG_KEY_SERVER_PORT     = "server.port"
	FLAG_KEY_CLIENT_USERNAME = "client.username"
	FLAG_KEY_CLIENT_PASSWORD = "client.password"
)

var (
	anumbers string
	mnumbers string
)

type auth struct {
	Username string
	Password string
}

type token struct {
	Token string `json:"token"`
}

type uuid struct {
	UUID string `json:"uuid"`
}

type numbers struct {
	Numbers string `json:"Numbers" form:"Numbers" query:"Numbers"`
}

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "gocele api client",
	Long:  "Simple cllient to interact with gocele api",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  "The version of the dispatch service.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add api call.",
	Long:  "Call the api to task a worker to add numbers.",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString(FLAG_KEY_SERVER_HOST)
		port := viper.GetString(FLAG_KEY_SERVER_PORT)
		username := viper.GetString(FLAG_KEY_CLIENT_USERNAME)
		password := viper.GetString(FLAG_KEY_CLIENT_PASSWORD)
		token := loginJSON(host, port, username, password)
		numberStr := anumbers

		// If we didn't get a token back, then error out
		if token == "" {
			log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
		}

		goAdd(host, port, token, numberStr)
	},
}

var mulCmd = &cobra.Command{
	Use:   "mul",
	Short: "Multiply api call.",
	Long:  "Call the api to task a worker to muliply numbers.",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString(FLAG_KEY_SERVER_HOST)
		port := viper.GetString(FLAG_KEY_SERVER_PORT)
		username := viper.GetString(FLAG_KEY_CLIENT_USERNAME)
		password := viper.GetString(FLAG_KEY_CLIENT_PASSWORD)
		token := loginJSON(host, port, username, password)
		numberStr := mnumbers

		// If we didn't get a token back, then error out
		if token == "" {
			log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
		}

		goMul(host, port, token, numberStr)
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print a JWT token.",
	Long:  "Print out a JWT token from a successful login",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString(FLAG_KEY_SERVER_HOST)
		port := viper.GetString(FLAG_KEY_SERVER_PORT)
		username := viper.GetString(FLAG_KEY_CLIENT_USERNAME)
		password := viper.GetString(FLAG_KEY_CLIENT_PASSWORD)
		token := loginJSON(host, port, username, password)

		if token == "" {
			log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
		}

		fmt.Printf("Your JWT Token is :[%s]\n", token)
	},
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup a task uuid.",
	Long:  "Looku a task uuid and return the result",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString(FLAG_KEY_SERVER_HOST)
		port := viper.GetString(FLAG_KEY_SERVER_PORT)
		username := viper.GetString(FLAG_KEY_CLIENT_USERNAME)
		password := viper.GetString(FLAG_KEY_CLIENT_PASSWORD)
		token := loginJSON(host, port, username, password)
		goLookup(host, port, token)
	},
}

/////// * API function calls below /////

func goLookup(host string, port string, tk string) {
	url := fmt.Sprintf("http://%s:%s/api/v1/tasks", host, port)

	id := uuid{viper.GetString("uuid")}
	jsonStr, _ := json.Marshal(id)

	auth := fmt.Sprintf("Bearer %s", tk)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

// This function will hopefully display a welcome message
// based on the authentication token provided in login

func goHello(host string, port string, tk string) {
	url := fmt.Sprintf("http://%s:%s/api/v1/hello", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func goAdd(host string, port string, tk string, nbrs string) {
	url := fmt.Sprintf("http://%s:%s/api/v1/add", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	data := numbers{nbrs}

	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func goMul(host string, port string, tk string, nbrs string) {
	url := fmt.Sprintf("http://%s:%s/api/v1/mul", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	data := numbers{nbrs}

	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

// This function will log you in via Json payload and return an auth token
// if successfull

func loginJSON(host string, port string, username string, password string) string {

	url := fmt.Sprintf("http://%s:%s/login", host, port)

	cred := auth{username, password}
	jsonStr, _ := json.Marshal(cred)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var t = new(token)
	err = json.Unmarshal(body, &t)

	if err != nil {
		log.Fatal(err)
	}
	return t.Token
}

func runCmd() {
	host := viper.GetString(FLAG_KEY_SERVER_HOST)
	port := viper.GetString(FLAG_KEY_SERVER_PORT)
	username := viper.GetString(FLAG_KEY_CLIENT_USERNAME)
	password := viper.GetString(FLAG_KEY_CLIENT_PASSWORD)
	token := loginJSON(host, port, username, password)
	goHello(host, port, token)
}

func init() {
	viper.SetConfigType("toml")
	viper.SetConfigName(SERVICE_NAME)
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", SERVICE_NAME))   // path to look for the config file in
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", SERVICE_NAME)) // call multiple times to add many search paths
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Fprintf(os.Stderr, "Fatal error config file: %s \n", err)
	}

	// Adding commands into the client
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(mulCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(lookupCmd)

	lookupFlags := lookupCmd.Flags()
	addFlags := addCmd.Flags()
	mulFlags := mulCmd.Flags()

	lookupFlags.String("uuid", "", "uuid of task you want to lookup.")
	viper.BindPFlag("uuid", lookupFlags.Lookup("uuid"))

	mulFlags.StringVar(&mnumbers, "i", "","integers to multiply")
	viper.BindPFlag("i", mulFlags.Lookup("i"))

	addFlags.StringVar(&anumbers,"i", "", "integers to add")
	viper.BindPFlag("i", addFlags.Lookup("i"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
