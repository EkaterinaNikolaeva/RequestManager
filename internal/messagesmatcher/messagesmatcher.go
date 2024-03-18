package messagesmatcher

import (
	"regexp"
	"strings"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
)

type MessagesMatcher struct {
	templateMatcher string
}

func (m MessagesMatcher) MatchMessage(message message.Message) bool {
	isMatch, _ := regexp.MatchString(m.templateMatcher, strings.ToLower(message.MessageText))
	return isMatch
}

func NewMessagesMatcher(templateMatcher string) MessagesMatcher {
	messagesMatcher := MessagesMatcher{templateMatcher: templateMatcher}
	return messagesMatcher
}
