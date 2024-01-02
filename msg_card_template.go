package feishu_bot_api

import (
	"encoding/json"
	"fmt"
)

var _ Message = (*cardMessageViaTemplate)(nil)

type cardMessageViaTemplate struct {
	id        string
	variables any
}

func (m cardMessageViaTemplate) Apply(body *MessageBody) error {
	rawVariables, err := json.Marshal(m.variables)
	if err != nil {
		return fmt.Errorf("card message(template): marshal variables: %w", err)
	}

	tpl := MessageBodyCardTemplate{
		Type: "template",
		Data: MessageBodyCardTemplateData{
			TemplateID:       m.id,
			TemplateVariable: (*json.RawMessage)(&rawVariables),
		},
	}
	rawTpl, err := json.Marshal(tpl)
	if err != nil {
		return fmt.Errorf("card message(template): marshal: %w", err)
	}

	body.MsgType = "interactive"
	body.Card = (*json.RawMessage)(&rawTpl)
	return nil
}
