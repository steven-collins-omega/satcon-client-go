package subscriptions

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryAddSubscription       = "addSubscription"
	AddSubscriptionVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","name":"{{js .Name}}","groups":[{{range $i,$e := .Groups}}{{if gt $i 0}},{{end}}"{{js $e}}"{{end}}],"channelUuid":"{{js .ChannelUUID}}","versionUuid":"{{js .VersionUUID}}"{{end}}`
)

type AddSubscriptionVariables struct {
	actions.GraphQLQuery
	OrgID       string
	Name        string
	Groups      []string
	ChannelUUID string
	VersionUUID string
}

func NewAddSubscriptionVariables(orgID, name, channelUuid, versionUuid string, groups []string) AddSubscriptionVariables {
	vars := AddSubscriptionVariables{
		OrgID:       orgID,
		Name:        name,
		Groups:      groups,
		ChannelUUID: channelUuid,
		VersionUUID: versionUuid,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryAddSubscription
	vars.Args = map[string]string{
		"orgId":       "String!",
		"name":        "String!",
		"groups":      "[String!]!",
		"channelUuid": "String!",
		"versionUuid": "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

// AddSubscriptionResponse for unmarshalling the response data
type AddSubscriptionResponse struct {
	Data *AddSubscriptionResponseData `json:"data,omitempty"`
}

// AddSubscriptionResponseData for unmarshalling response details
type AddSubscriptionResponseData struct {
	Details *AddSubscriptionResponseDataDetails `json:"addSubscription,omitempty"`
}

// AddSubscriptionResponseDataDetails for unmarshalling response uuid
type AddSubscriptionResponseDataDetails struct {
	UUID string `json:"uuid,omitempty"`
}

// AddSubscription creates a new subscription for valid channel, version, and group(s)
func (c *Client) AddSubscription(orgID, name, channelUuid, versionUuid string, groups []string, token string) (*AddSubscriptionResponseDataDetails, error) {
	var response AddSubscriptionResponse

	vars := NewAddSubscriptionVariables(orgID, name, channelUuid, versionUuid, groups)

	err := c.DoQuery(AddSubscriptionVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, nil
	}

	return nil, err
}
