package main

import (
	"fmt"
	"testing"
)

func TestReadMessage(t *testing.T) {
	msg := "{\"id\":\"947c90a5-bfdb-4377-9573-1bf48d59046a\",\"projectId\":\"tremaps4\",\"status\":\"SUCCESS\",\"source\":{\"repoSource\":{\"projectId\":\"tremaps4\",\"repoName\":\"github_trealtamira_tremaps-task\",\"branchName\":\"master\"}},\"steps\":[{\"name\":\"gcr.io/cloud-builders/docker\",\"args\":[\"build\",\"--tag=gcr.io/tremaps4/task\",\".\"],\"timing\":{\"startTime\":\"2019-09-19T10:17:37.900006022Z\",\"endTime\":\"2019-09-19T10:18:12.368687262Z\"},\"pullTiming\":{\"startTime\":\"2019-09-19T10:17:37.900006022Z\",\"endTime\":\"2019-09-19T10:17:37.967306313Z\"},\"status\":\"SUCCESS\"}],\"results\":{\"buildStepImages\":[\"sha256:f7e685514163b3ff8fffb47c408b8ec55c83f005a97e8368b4e809ebbe41ff04\"],\"buildStepOutputs\":[]},\"createTime\":\"2019-09-19T10:17:25.314271826Z\",\"startTime\":\"2019-09-19T10:17:27.124787545Z\",\"finishTime\":\"2019-09-19T10:18:13.791250Z\",\"timeout\":\"600s\",\"logsBucket\":\"gs://244961050376.cloudbuild-logs.googleusercontent.com\",\"sourceProvenance\":{\"resolvedRepoSource\":{\"projectId\":\"tremaps4\",\"repoName\":\"github_trealtamira_tremaps-task\",\"commitSha\":\"20b1c4fee56caa197e7c498e66ebc709e4103499\"}},\"buildTriggerId\":\"a128b807-995e-4021-a2d6-a22a8aa5191d\",\"options\":{\"substitutionOption\":\"ALLOW_LOOSE\",\"logging\":\"LEGACY\"},\"logUrl\":\"https://console.cloud.google.com/gcr/builds/947c90a5-bfdb-4377-9573-1bf48d59046a?project=244961050376\",\"tags\":[\"trigger-a128b807-995e-4021-a2d6-a22a8aa5191d\"],\"timing\":{\"BUILD\":{\"startTime\":\"2019-09-19T10:17:37.899965985Z\",\"endTime\":\"2019-09-19T10:18:12.447988697Z\"},\"FETCHSOURCE\":{\"startTime\":\"2019-09-19T10:17:29.902813257Z\",\"endTime\":\"2019-09-19T10:17:37.835708447Z\"}}}"
	//msg = "{\"id\":\"947c90a5-bfdb-4377-9573-1bf48d59046a\"}"

	chat, _ := readMessage(msg)
	fmt.Printf("LogUrl: %v\n", chat.LogURL)
	fmt.Printf("Project: %v\n", chat.ProjectId)
	fmt.Printf("Status: %v\n", chat.Status)
	fmt.Printf("Branch: %v\n", chat.Branch)
	fmt.Printf("ProjectID: %v\n", chat.Project)
	fmt.Printf("Repo: %v\n", chat.Repo)
}

func TestBuildChat(t *testing.T) {
	cmess := []BuildData{{LogURL: "url://sito", ProjectId: "PID", StartTime: "2019-09-19T10:18:12.368687262Z", Project: "PROGETTO", Repo: "RE PO", FinishTime: "2019-09-19T10:18:12.368687262Z", Branch: "RA MO", Status: "FAILED"},
		{LogURL: "url://sito", ProjectId: "PID", StartTime: "2019-09-19T10:18:12.368687262Z", Project: "PROGETTO", Repo: "RE PO", FinishTime: "2019-09-19T10:18:12.368687262Z", Branch: "RA MO", Status: "SUCCESS"},
		{LogURL: "url://sito", ProjectId: "PID", StartTime: "2019-09-19T10:18:12.368687262Z", Project: "PROGETTO", Repo: "RE PO", FinishTime: "2019-09-19T10:18:12.368687262Z", Branch: "RA MO", Status: "WAITING"},
	}
	for _, c := range cmess {
		message, _ := buildChat(c)
		fmt.Println(message)
	}
}

func TestSendChat(t *testing.T) {
	jsonchat := "{\"cards\":[ { \"header\": { \"title\": \"Cloud Builder TEST\" }, \"sections\": [ { \"widgets\": [ { \"textParagraph\": { \"text\": \"TEST\"}}]}]}]}"
	err := sendChatJSON(jsonchat)
	if err != nil {
		t.Errorf("something went wrong %v", err)
	}
}

func TestAllFlow(t *testing.T) {
	msg := "{\"id\":\"947c90a5-bfdb-4377-9573-1bf48d59046a\",\"projectId\":\"tremaps4\",\"status\":\"SUCCESS\",\"source\":{\"repoSource\":{\"projectId\":\"tremaps4\",\"repoName\":\"github_trealtamira_tremaps-task\",\"branchName\":\"master\"}},\"steps\":[{\"name\":\"gcr.io/cloud-builders/docker\",\"args\":[\"build\",\"--tag=gcr.io/tremaps4/task\",\".\"],\"timing\":{\"startTime\":\"2019-09-19T10:17:37.900006022Z\",\"endTime\":\"2019-09-19T10:18:12.368687262Z\"},\"pullTiming\":{\"startTime\":\"2019-09-19T10:17:37.900006022Z\",\"endTime\":\"2019-09-19T10:17:37.967306313Z\"},\"status\":\"SUCCESS\"}],\"results\":{\"buildStepImages\":[\"sha256:f7e685514163b3ff8fffb47c408b8ec55c83f005a97e8368b4e809ebbe41ff04\"],\"buildStepOutputs\":[]},\"createTime\":\"2019-09-19T10:17:25.314271826Z\",\"startTime\":\"2019-09-19T10:17:27.124787545Z\",\"finishTime\":\"2019-09-19T10:18:13.791250Z\",\"timeout\":\"600s\",\"logsBucket\":\"gs://244961050376.cloudbuild-logs.googleusercontent.com\",\"sourceProvenance\":{\"resolvedRepoSource\":{\"projectId\":\"tremaps4\",\"repoName\":\"github_trealtamira_tremaps-task\",\"commitSha\":\"20b1c4fee56caa197e7c498e66ebc709e4103499\"}},\"buildTriggerId\":\"a128b807-995e-4021-a2d6-a22a8aa5191d\",\"options\":{\"substitutionOption\":\"ALLOW_LOOSE\",\"logging\":\"LEGACY\"},\"logUrl\":\"https://console.cloud.google.com/gcr/builds/947c90a5-bfdb-4377-9573-1bf48d59046a?project=244961050376\",\"tags\":[\"trigger-a128b807-995e-4021-a2d6-a22a8aa5191d\"],\"timing\":{\"BUILD\":{\"startTime\":\"2019-09-19T10:17:37.899965985Z\",\"endTime\":\"2019-09-19T10:18:12.447988697Z\"},\"FETCHSOURCE\":{\"startTime\":\"2019-09-19T10:17:29.902813257Z\",\"endTime\":\"2019-09-19T10:17:37.835708447Z\"}}}"
	chat, _ := readMessage(msg)
	message, _ := buildChat(chat)
	fmt.Println(message)
	sendChatJSON(message)

}
