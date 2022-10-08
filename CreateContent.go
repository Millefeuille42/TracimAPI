package TracimAPI

import (
	"encoding/json"
	"fmt"
)

type createContentData struct {
	ParentId    interface{} `json:"parent_id"`
	ContentType string      `json:"content_type"`
	Label       string      `json:"label"`
}

func (a *Api) CreateContent(workspaceId, tType, name string) error {
	contentData := createContentData{
		ParentId:    nil,
		ContentType: tType,
		Label:       name,
	}
	data, err := json.Marshal(contentData)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/workspaces/%s/contents", workspaceId)
	_, err = a.Request("POST", endpoint, data)
	if err != nil {
		return err
	}
	return nil
}
