package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	JIRA_UPDATE_STATUS_URL  = "http://almsmart1.demos.hclets.com/jira/rest/api/2/issue/" + ISSUE_PLACEHOLDER + "/transitions?expand=transitions.fields"
	JIRA_UPDATE_COMMENT_URL = "http://almsmart1.demos.hclets.com/jira/rest/api/2/issue/" + ISSUE_PLACEHOLDER + "/comment"
	ZAPI_TEST_UPDATE_URL    = "http://almsmart1.demos.hclets.com/jira/rest/zapi/latest/execution/" + STEP_ID_PLACEHOLDER + "/execute"
	STEP_ID_PLACEHOLDER     = "${step-id}"
	ISSUE_PLACEHOLDER       = "${issue-id}"
	USERNAME                = "Admin"
	PASSWORD                = "ALMcoe@123"
)

func UpdateJiraIssues(tcs []TestCase) error {
	for _, tc := range tcs {
		/*url := prepareTransitionJiraURL(tc.TestID)
		fmt.Printf("URL: %s\n", url)
		err := updateIssueAttributes(tc, url, "transition")
		if err != nil {
			fmt.Printf("Error Changing status for issue %s, Error %s\n", tc.TestID, err)
		} else {*/

		stepUrl := prepareStepUpdateUrl(tc.TestStepID)
		err := updateIssueAttributes(tc, stepUrl, "step")
		if err != nil {
			fmt.Printf("Error Changing status for test step %s, Error %s\n", tc.TestStepID, err)
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
	} else if attrType == "comment" {
		body = prepareCommentReqBody(tc)
	} else {
		body = prepareStepIdUpdateReqBody(tc)
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.Encode(body)
	//fmt.Printf("body | %s\n", buffer.String())
	httpClient := &http.Client{}
	method := http.MethodPost
	if attrType == "step" {
		method = http.MethodPut
	}
	request, _ := http.NewRequest(method, url, &buffer)

	if attrType != "step" {
		cookies := getCookies(url)
		if len(cookies) > 0 {
			for _, co := range cookies {
				request.AddCookie(co)
				if attrType == "step" {
					fmt.Printf("Cookies for step update Name: %s \n", co.Name)
				}
			}
		} else {
			fmt.Printf("No Cookies found\n")
		}
	}
	request.Header.Add("Content-Type", "application/json")
	//request.Header.Add("Authorization", "Basic QWRtaW46QUxNY29lQDEyMw==")
	//request.Header.Add("User-Agent", "PostmanRuntime/7.15.0")
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
	//defer r.Body.Close()
	if err != nil {
		fmt.Printf("Error process request to url %s Error: %s\n", url, err)
		os.Exit(128)
	}
	if r.StatusCode == http.StatusOK || r.StatusCode == http.StatusNoContent || r.StatusCode == http.StatusCreated {
		if attrType == "step" && len(r.Cookies()) < 2 {
			fmt.Printf("Failure: Failed to update the test step for issue %s\n", tc.TestID)
		} else {
			fmt.Printf("Success: Issue %s Updated Successfully | Status Code %d \n", tc.TestID, r.StatusCode)
		}
	} else {
		fmt.Printf("Success: Status Code %d \n", r.StatusCode)
		//b, err := ioutil.ReadAll(r.Body)
		cs := r.Cookies()
		for _, co := range cs {
			fmt.Printf("Cookie Name %s\n", co.Name)
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
	request, _ := http.NewRequest(http.MethodHead, url, nil)
	request.SetBasicAuth(USERNAME, PASSWORD)
	r, err := httpClient.Do(request)
	if err != nil {
		return []*http.Cookie{}
	}
	return r.Cookies()
}

func prepareStepIdUpdateReqBody(tc TestCase) Status {
	reqBody := Status{}
	s := strings.ToLower(tc.Status)
	s = strings.TrimSpace(s)
	if s == "ok" {
		reqBody.Status = "1"
	} else {
		reqBody.Status = "2"
	}
	return reqBody
}

func prepareStepUpdateUrl(testStepId string) string {
	url := strings.Replace(ZAPI_TEST_UPDATE_URL, STEP_ID_PLACEHOLDER, testStepId, -1)
	return url
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
	req.Body = "Issue" + tc.TestID + "updated in " + tc.ResTime + " msec."
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

type Status struct {
	Status string `json:"status"`
}
