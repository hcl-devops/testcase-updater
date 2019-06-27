package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	JIRA_UPDATE_STATUS_URL  = "http://almsmart1.demos.hclets.com/jira/rest/api/2/issue/" + ISSUE_PLACEHOLDER + "/transitions?expand=transitions.fields"
	JIRA_UPDATE_COMMENT_URL = "http://almsmart1.demos.hclets.com/jira/rest/api/2/issue/" + ISSUE_PLACEHOLDER + "/comment"
	ISSUE_PLACEHOLDER       = "${issue-id}"
	USERNAME                = "Admin"
	PASSWORD                = "ALMcoe@123"
)

func UpdateJiraIssues(tcs []TestCase) error {
	for _, tc := range tcs {
		url := prepareTransitionJiraURL(tc.TestID)
		//fmt.Printf("URL: %s\n", url)
		err := updateIssueAttributes(tc, url, "transition")
		if err != nil {
			fmt.Printf("Error Changing status for issue %s, Error %s\n", tc.TestID, err)
		} else {
			urlComment := prepareCommentJiraURL(tc.TestID)
			err = updateIssueAttributes(tc, urlComment, "comment")
			if err != nil {
				fmt.Printf("Error Adding comment for issue %s, Error %s\n", tc.TestID, err)
			}
		}
	}
	return nil
}

func updateIssueAttributes(tc TestCase, url, attrType string) error {
	var body interface{}
	if attrType == "transition" {
		body = prepareTransitionReqBody(tc, tc.Status)
	} else {
		body = prepareCommentReqBody(tc)
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.Encode(body)
	//fmt.Printf("body | %s\n", buffer.String())
	httpClient := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, url, &buffer)

	cookies := getCookies(url)
	if len(cookies) > 0 {
		for _, co := range cookies {
			request.AddCookie(co)
		}
	} else {
		fmt.Printf("No Cookies found\n")
	}
	request.Header.Add("Content-Type", "application/json")
	//request.Header.Add("Authorization", "Basic QWRtaW46QUxNY29lQDEyMw==")
	request.Header.Add("User-Agent", "PostmanRuntime/7.15.0")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Postman-Token", "20403faf-d0e4-4c43-b37e-5bc8e35c92cb,b5c01a22-17a2-4f93-835a-9ec9f15eee99")
	request.Header.Add("Host", "almsmart1.demos.hclets.com")
	//request.Header.Add("accept-encoding", "gzip, deflate")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("cache-control", "no-cache")
	//request.Header.Set("Authorization", "Basic QWRtaW46QUxNY29lQDEyMw==")
	request.SetBasicAuth(USERNAME, PASSWORD)
	//fmt.Printf("Auth: %s\n", request.Header.Get("Authorization"))
	r, err := httpClient.Do(request)
	defer r.Body.Close()

	if r.StatusCode == http.StatusOK || r.StatusCode == http.StatusNoContent || r.StatusCode == http.StatusCreated {
		fmt.Printf("Success: Issue %s Updated Successfully | Status Code %d \n", tc.TestID, r.StatusCode)
	} else {
		fmt.Printf("Success: Status Code %d \n", r.StatusCode)
		//b, err := ioutil.ReadAll(r.Body)
		cs := r.Cookies()
		for _, co := range cs {
			fmt.Printf("Cookie Name %s, Path %s, Raw %s\n", co.Name, co.Path, co.Raw)
		}
		if err != nil {
			fmt.Printf("Error reading response\n")
			return err
		} else {
			// do nothing
		}
	}
	if err != nil {
		fmt.Printf("Error updating Jira issue %s, error: %s\n", tc.TestID, err)
		return err
	} else {
		//fmt.Printf("Jira issue %s, udpated succeessfully\n", tc.TestID)
	}
	return nil

}

func getCookies(url string) []*http.Cookie {
	httpClient := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.SetBasicAuth(USERNAME, PASSWORD)
	r, err := httpClient.Do(request)
	if err != nil {
		return []*http.Cookie{}
	}
	return r.Cookies()
}

func prepareTransitionReqBody(tc TestCase, status string) UpdateReqBody {
	reqBody := UpdateReqBody{}
	transition := Transition{}
	s := strings.ToLower(status)
	s = strings.TrimSpace(s)
	fmt.Printf("Status of test case %s is%s.\n", tc.TestID, s)
	if s == "ok" {
		transition.Id = "31"
	} else {
		transition.Id = "21"
	}
	reqBody.Transition = transition
	return reqBody
}

func prepareCommentReqBody(tc TestCase) Comment {
	req := Comment{}
	req.Body = "Issue" + tc.TestID + "updated in 0.01 sec"
	return req
}

func prepareTransitionJiraURL(testId string) string {
	url := strings.Replace(JIRA_UPDATE_STATUS_URL, ISSUE_PLACEHOLDER, testId, -1)
	return url
}

func prepareCommentJiraURL(testId string) string {
	url := strings.Replace(JIRA_UPDATE_COMMENT_URL, ISSUE_PLACEHOLDER, testId, -1)
	return url
}

type UpdateReqBody struct {
	//Update     Update     `json:"update"`
	Transition Transition `json:"transition"`
}

type Comment struct {
	Body string `json:"body"`
}

type Transition struct {
	Id string `json:"id"`
}
