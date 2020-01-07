package alerts

import "fmt"

// ListPluginsConditions returns alert conditions for New Relic plugins for an account.
func (alerts *Alerts) ListPluginsConditions(policyID int) ([]*PluginCondition, error) {
	response := pluginsConditionsResponse{}
	conditions := []*PluginCondition{}
	queryParams := listPluginsConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_plugins_conditions.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		for _, c := range response.PluginsConditions {
			c.PolicyID = policyID
		}

		conditions = append(conditions, response.PluginsConditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// GetPluginCondition gets information about an alert condition for a plugin
// given a policy ID and plugin ID.
func (alerts *Alerts) GetPluginCondition(policyID int, id int) (*PluginCondition, error) {
	conditions, err := alerts.ListPluginsConditions(policyID)

	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return condition, nil
		}
	}

	return nil, fmt.Errorf("no condition found for policy %d and condition ID %d", policyID, id)
}

type listPluginsConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type pluginsConditionsResponse struct {
	PluginsConditions []*PluginCondition `json:"plugins_conditions,omitempty"`
}
