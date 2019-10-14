package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", PushToChat)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// PushToChat receives and processes a Pub/Sub build message to Google Chat.
func PushToChat(w http.ResponseWriter, r *http.Request) {
	var m PubSubMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("iotuil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	cm, err := readMessage(string(m.Message.Data))
	chat, err := buildChat(cm)
	if chat != "" {
		err = sendChatJSON(chat)
	}
}

type BuildData struct {
	LogURL     string
	ProjectId  string
	StartTime  string
	FinishTime string
	Status     string
	Branch     string
	Project    string
	Repo       string
}

func readMessage(data string) (BuildData, error) {
	var msgMap map[string]interface{}
	log.Printf("readMessage: %s\n\n", data)
	json.Unmarshal([]byte(data), &msgMap)
	chat := BuildData{}
	for k, v := range msgMap {
		switch k {
		case "logUrl":
			chat.LogURL = v.(string)
			break
		case "projectId":
			chat.ProjectId = v.(string)
			break
		case "status":
			chat.Status = v.(string)
			break
		case "startTime":
			chat.StartTime = v.(string)
			break
		case "finishTime":
			chat.FinishTime = v.(string)
			break
		case "source":
			source := v.(map[string]interface{})
			repo := source["repoSource"].(map[string]interface{})
			if repo != nil {
				chat.Branch = repo["branchName"].(string)
				chat.Repo = repo["repoName"].(string)
				chat.Project = repo["projectId"].(string)
			}
			break
		}
	}
	return chat, nil
}

type ChatMessage struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	Sections []Section `json:"sections,omitempty"`
	Header   Header    `json:"header,omitempty"`
}

type Header struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type Section struct {
	Header  string                   `json:"header,omitempty"`
	Widgets []map[string]interface{} `json:"widgets,omitempty"`
}

type Paragraph struct {
	Text string `json:"text,omitempty"`
}

type Image struct {
	ImageURL string `json:"imageUrl,omitempty"`
}

type KeyValue struct {
	TopLabel string `json:"topLabel,omitempty"`
	Content  string `json:"content,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

func buildChat(c BuildData) (string, error) {
	if c.Status == "WORKING" || c.Status == "QUEUED" {
		log.Printf("Status is working or queued.. sending nothing.")
		return "", nil
	}
	section := Section{}
	section.Widgets = make([]map[string]interface{}, 0)

	//{"id":"59d19f38-e1cc-4213-b315-48e686fb974d","projectId":"tremaps4","status":"SUCCESS","source":{"repoSource":{"projectId":"tremaps4","repoName":"github_trealtamira_tremaps-task","branchName":"master"}},"steps":[{"name":"gcr.io/cloud-builders/docker","args":["build","--tag=gcr.io/tremaps4/task","."],"timing":{"startTime":"2019-09-27T09:42:37.249119350Z","endTime":"2019-09-27T09:43:00.147304522Z"},"pullTiming":{"startTime":"2019-09-27T09:42:37.249119350Z","endTime":"2019-09-27T09:42:37.309943561Z"},"status":"SUCCESS"}],"results":{"buildStepImages":["sha256:f7e685514163b3ff8fffb47c408b8ec55c83f005a97e8368b4e809ebbe41ff04"]},"createTime":"2019-09-27T09:42:30.758288385Z","startTime":"2019-09-27T09:42:32.378117093Z","finishTime":"2019-09-27T09:43:01.944340Z","timeout":"600s","logsBucket":"gs://244961050376.cloudbuild-logs.googleusercontent.com","sourceProvenance":{"resolvedRepoSource":{"projectId":"tremaps4","repoName":"github_trealtamira_tremaps-task","commitSha":"20b1c4fee56caa197e7c498e66ebc709e4103499"}},"buildTriggerId":"a128b807-995e-4021-a2d6-a22a8aa5191d","options":{"substitutionOption":"ALLOW_LOOSE","logging":"LEGACY"},"logUrl":"https://console.cloud.google.com/gcr/builds/59d19f38-e1cc-4213-b315-48e686fb974d?project=244961050376","tags":["trigger-a128b807-995e-4021-a2d6-a22a8aa5191d"],"timing":{"BUILD":{"startTime":"2019-09-27T09:42:37.249083264Z","endTime":"2019-09-27T09:43:00.237587708Z"},"FETCHSOURCE":{"startTime":"2019-09-27T09:42:33.583442569Z","endTime":"2019-09-27T09:42:37.179067274Z"}}}

	//TODO add red when status is not success
	if c.Status != "" {
		var status string
		if c.Status == "SUCCESS" {
			status = fmt.Sprintf("<font color=\"#00ff00\"><b>%s</b></font>", c.Status)
		} else if c.Status == "FAILED" {
			status = fmt.Sprintf("<font color=\"#ff0000\"><b>%s</b></font>", c.Status)
		} else {
			status = fmt.Sprintf("<font color=\"#ffff00\"><b>%s</b></font>", c.Status)
		}
		kv := KeyValue{TopLabel: "Status", Content: status, Icon: "MEMBERSHIP"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	if c.ProjectId != "" {
		kv := KeyValue{TopLabel: "ProjectId", Content: c.ProjectId, Icon: "STAR"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	if c.StartTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05.000000000Z", c.StartTime)
		st := t.Format("2006-01-02 15:04:05 UTC")
		if err != nil {
			st = c.StartTime
		}
		kv := KeyValue{TopLabel: "Start Time", Content: st, Icon: "CLOCK"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	if c.FinishTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05.000000Z", c.FinishTime)
		ft := t.Format("2006-01-02 15:04:05 UTC")
		if err != nil {
			ft = c.FinishTime
		}
		kv := KeyValue{TopLabel: "Finish Time", Content: ft, Icon: "CLOCK"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	if c.Repo != "" {
		kv := KeyValue{TopLabel: "Repo", Content: c.Repo, Icon: "DESCRIPTION"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	if c.Branch != "" {
		kv := KeyValue{TopLabel: "Branch", Content: c.Branch, Icon: "DESCRIPTION"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	if c.LogURL != "" {
		url := fmt.Sprintf("<a href=\"%s\">Build log</a>", c.LogURL)
		kv := KeyValue{TopLabel: "Build URL", Content: url, Icon: "BOOKMARK"}
		section.Widgets = append(section.Widgets, map[string]interface{}{"keyValue": kv})
	}
	fmt.Println(section.Widgets)

	card := Card{}
	card.Header = Header{Title: "Build Completed"}
	card.Sections = make([]Section, 0)
	card.Sections = append(card.Sections, section)

	chatMessage := ChatMessage{}
	chatMessage.Cards = make([]Card, 0)
	chatMessage.Cards = append(chatMessage.Cards, card)

	msg, err := json.Marshal(chatMessage)
	if err != nil {
		fmt.Printf("Encoding failed: %v", err)
	}
	return string(msg), err
}

func sendChatJSON(message string) error {
	webhook := os.Getenv("WEBHOOK")
	fmt.Printf("Webhook is: %s\n", webhook)
	resp, err := http.Post(webhook, "application/json", strings.NewReader(message))
	if err != nil {
		log.Printf("post failed: %v\n\n", err)
	}
	fmt.Printf("Response from chat hook: %v", resp)
	return err
}
