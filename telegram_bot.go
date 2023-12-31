package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/yaml.v2"
)

var triggers []Trigger

type Trigger struct {
	Key    string      `yaml:"key"`
	Values interface{} `yaml:"values"`
}

// Update struct:
//
//	{
//	    "update_id": 123,
//	    "message": {
//	        "message_id": 123,
//	        "from": {
//	            "id": 123,
//	            "is_bot": false,
//	            "first_name": "John",
//	            "username": "john_doe",
//	            "language_code": "en"
//	        },
//	        "chat": {
//	            "id": triggerName,
//	            "first_name": "John",
//	            "username": "john_doe",
//	            "type": "private"
//	        },
//	        "date": 1703257400,
//	        "text": "test"
//	    }
//	}
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

func initializeTriggers() error {
	filename, _ := filepath.Abs("./triggers.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &triggers)
	return err
}

func main() {
	loadEnvVariables()

	err := initializeTriggers()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", callHandler)

	fmt.Printf("Server is running on port %s\n", os.Getenv("PORT"))
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func loadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	if os.Getenv("TOKEN") == "" {
		log.Fatalf("Missing TOKEN environment variable")
	}
	if os.Getenv("PORT") == "" {
		log.Fatalf("Missing PORT environment variable")
	}
}

func callHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// logRequest(request) // Debug

	var update Update
	err := json.NewDecoder(request.Body).Decode(&update)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	p := bluemonday.UGCPolicy()
	sanitizedMessageText := p.Sanitize(update.Message.Text)
	responseText := computeResponse(sanitizedMessageText)

	if responseText != "" {
		var telegramResponseBody, telegramError = sendTextToChat(update.Message.Chat.Id, responseText)
		if telegramError != nil {
			log.Printf("Error %s from Telegram; reponse body is %s", telegramError.Error(), telegramResponseBody)
		} else {
			log.Printf("Message successfuly sent to chat id %d", update.Message.Chat.Id)
		}
	}
}

func computeResponse(inputMessage string) string {
	lowerMessage := strings.ToLower(inputMessage)
	for _, trigger := range triggers {
		lowerKey := strings.ToLower(trigger.Key)
		index := strings.Index(lowerMessage, lowerKey)
		if index > -1 { // The trigger matches
			switch values := trigger.Values.(type) {
			case string:
				return formatResponse(values, inputMessage, trigger.Key, index)
			case []interface{}:
				return formatResponse(getRandomFromArray(values), inputMessage, trigger.Key, index)
			}
		}
	}
	return ""
}

func formatResponse(response string, inputMessage string, triggerKey string, index int) string {
	if !strings.Contains(response, "__input__") {
		return response
	} else {
		input := inputMessage[index+len(triggerKey):]
		return strings.Replace(response, "__input__", strings.TrimSpace(input), -1)
	}
}

func getRandomFromArray(array []interface{}) string {
	randomIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(array))
	if str, ok := array[randomIndex].(string); ok {
		return str
	}
	return ""
}

func sendTextToChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	response, err := http.PostForm(
		apiUrl("sendMessage"),
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("Error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, readError = io.ReadAll(response.Body)
	if readError != nil {
		log.Printf("Error parsing Telegram response: %s", readError.Error())
		return "", readError
	}
	bodyString := string(bodyBytes)
	log.Printf("Telegram Response: %s", bodyString)

	return bodyString, nil
}

func apiUrl(endpoint string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", os.Getenv("TOKEN"), endpoint)
}

func logRequest(request *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	reqBytes := buf.Bytes()

	// Reset the request body position to the beginning
	request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))

	reqString := string(reqBytes)
	fmt.Printf("%s\n", reqString)
}
