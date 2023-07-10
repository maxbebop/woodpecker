package jiraclient

import (
	"fmt"

	jira "gopkg.in/andygrunwald/go-jira.v1"
)

type JiraIssues interface {
	Search(jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error)
	AddComment(issueID string, comment *jira.Comment) (*jira.Comment, *jira.Response, error)
}

type Logger interface {
	Printf(s string, params ...any)
	Err(msg interface{}, keyvals ...interface{}) error
}

type Client struct {
	api JiraIssues
	log Logger
}

func New(baseUrl, username, token string, log Logger) (*Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		return nil, fmt.Errorf("new jira client")
	}

	return &Client{api: jiraClient.Issue, log: log}, nil
}

/*
func (c *Client) GetIssues(jql string) ([]jira.Issue, error) {

	lastIssue := 0
	var result []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: maxResults,
			StartAt:    lastIssue,
		}
		issues, resp, err := c.api.Issue.Search(jql, opt)
		if err != nil {
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			result = make([]jira.Issue, 0, total)
		}

		result = append(result, issues...)
		lastIssue = resp.StartAt + len(issues)

		if lastIssue >= total {
			break
		}
	}

	for _, i := range result {
		fmt.Printf("%s (%s/%s): %+v -> %s\n", i.Key, i.Fields.Type.Name, i.Fields.Priority.Name, i.Fields.Summary, i.Fields.Status.Name)
		fmt.Printf("Assignee : %v\n", i.Fields.Assignee.DisplayName)
		fmt.Printf("Reporter: %v\n", i.Fields.Reporter.DisplayName)
	}

	return result, nil
}

func (c *Client) AddCommentToIssue(issueId string, comment string) error {
	_, _, err := c.api.Issue.AddComment(issueId, &jira.Comment{
		Body: comment,
	})

	return err
}
*/
