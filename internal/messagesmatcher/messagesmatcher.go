package messagesmatcher

import (
	"regexp"
	"strings"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/message"
)

type MessagesMatcher struct {
}

var validMsg = regexp.MustCompile(`.*jira!.*`)

func (m MessagesMatcher) MatchMessage(message message.Message) bool {
	return validMsg.MatchString(strings.ToLower(message.Message))
}
