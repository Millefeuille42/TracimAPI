package session

import (
	"encoding/json"
	"fmt"
	"time"
)

type createContentRespData struct {
	ContentNamespace string      `json:"content_namespace"`
	ParentId         interface{} `json:"parent_id"`
	Modified         time.Time   `json:"modified"`
	IsEditable       bool        `json:"is_editable"`
	FileExtension    string      `json:"file_extension"`
	Created          time.Time   `json:"created"`
	LastModifier     struct {
		HasCover   bool   `json:"has_cover"`
		Username   string `json:"username"`
		PublicName string `json:"public_name"`
		UserId     int    `json:"user_id"`
		HasAvatar  bool   `json:"has_avatar"`
	} `json:"last_modifier"`
	Slug                string      `json:"slug"`
	WorkspaceId         int         `json:"workspace_id"`
	ShowInUi            bool        `json:"show_in_ui"`
	Status              string      `json:"status"`
	CurrentRevisionType string      `json:"current_revision_type"`
	ActivesShares       int         `json:"actives_shares"`
	IsArchived          bool        `json:"is_archived"`
	Label               string      `json:"label"`
	CurrentRevisionId   int         `json:"current_revision_id"`
	AssigneeId          interface{} `json:"assignee_id"`
	ContentType         string      `json:"content_type"`
	ContentId           int         `json:"content_id"`
	Filename            string      `json:"filename"`
	IsTemplate          bool        `json:"is_template"`
	IsDeleted           bool        `json:"is_deleted"`
	SubContentTypes     []string    `json:"sub_content_types"`
	Author              struct {
		HasCover   bool   `json:"has_cover"`
		Username   string `json:"username"`
		PublicName string `json:"public_name"`
		UserId     int    `json:"user_id"`
		HasAvatar  bool   `json:"has_avatar"`
	} `json:"author"`
}

type createContentData struct {
	ParentId    interface{} `json:"parent_id"`
	ContentType string      `json:"content_type"`
	Label       string      `json:"label"`
}

func (a *Session) CreateContent(workspaceId, tType, name string) (string, error) {
	contentData := createContentData{
		ParentId:    nil,
		ContentType: tType,
		Label:       name,
	}
	data, err := json.Marshal(contentData)
	if err != nil {
		return "", err
	}
	endpoint := fmt.Sprintf("/workspaces/%s/contents", workspaceId)
	response, err := a.Request("POST", endpoint, data)
	if err != nil {
		return "", err
	}
	respData := createContentRespData{}
	err = json.Unmarshal(response.DataBytes, &respData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", respData.ContentId), nil
}
