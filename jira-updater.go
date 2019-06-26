package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	JIRA_UPDATE_STATUS_COMMENT_URL = "http://almsmart1.demos.hclets.com/jira/rest/api/latest/issue/" + ISSUE_PLACEHOLDER + "/transitions?expand=transitions.fields"
	ISSUE_PLACEHOLDER              = "${issue-id}"
	USERNAME                       = "Admin"
	PASSWORD                       = "ALMcoe@123"
)

func UpdateJiraIssues(tcs []TestCase) error {
	for _, tc := range tcs {
		url := prepareJiraURL(tc.TestID)
		fmt.Printf("Url for case %s is %s\n", tc.TestID, url)
		status := tc.Status
		body := prepareReqBody(tc, status)

		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.Encode(body)
		fmt.Printf("json %s \n", buffer.String())

		httpClient := &http.Client{}
		request, _ := http.NewRequest(http.MethodPost, url, &buffer)
		request.Header.Set("Content-Type", "application/json")
		request.SetBasicAuth(USERNAME, PASSWORD)

		r, err := httpClient.Do(request)
		if r.StatusCode != http.StatusOK || r.StatusCode != http.StatusNoContent {
			fmt.Printf("Status Code %d Status Mesage %s \n", r.StatusCode, r.Status)
			//fmt.Printf("%s\n", r.Body)
		}
		if err != nil {
			fmt.Printf("Error updating Jira issue %s, error: %s\n", tc.TestID, err)
		} else {
			//fmt.Printf("Jira issue %s, udpated succeessfully\n", tc.TestID)
		}
	}
	return nil
}

func prepareReqBody(tc TestCase, status string) UpdateReqBody {
	reqBody := UpdateReqBody{}
	id := tc.TestID
	transition := Transition{}
	s := strings.ToLower(status)
	fmt.Printf("Status of test case %s is %s\n", tc.TestID, s)
	if s == "ok" {
		transition.Id = "31"
	} else {
		transition.Id = "21"
	}
	u := Update{}
	var cs []Comment
	c := Comment{}
	c.Add = Add{Body: "test case " + id + "executed in 0.01 sec"}
	cs = append(cs, c)
	u.Comment = cs
	reqBody.Transition = transition
	reqBody.Update = u
	//fmt.Printf("body: %v\n", reqBody)
	return reqBody
}

func prepareJiraURL(testId string) string {
	url := strings.Replace(JIRA_UPDATE_STATUS_COMMENT_URL, ISSUE_PLACEHOLDER, testId, -1)
	return url
}

type UpdateReqBody struct {
	Update     Update     `json:"update"`
	Transition Transition `json:"transition"`
}

type Update struct {
	Comment []Comment `json:"comment"`
}

type Comment struct {
	Add Add `json:"add"`
}

type Add struct {
	Body string `json:"body"`
}

type Transition struct {
	Id string `json:"id"`
}
