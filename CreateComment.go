package TracimAPI

import (
	"encoding/json"
	"fmt"
)

type createCommentData struct {
	RawContent       string `json:"raw_content"`
	ContentNamespace string `json:"content_namespace"`
}

func (a *Api) CreateComment(workspaceId, contentId, text string) error {
	commentData := createCommentData{
		RawContent:       text,
		ContentNamespace: "content",
	}
	data, err := json.Marshal(commentData)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/workspaces/%s/contents/%s/comments", workspaceId, contentId)
	_, err = a.Request("POST", endpoint, data)
	if err != nil {
		return err
	}
	return nil
}
