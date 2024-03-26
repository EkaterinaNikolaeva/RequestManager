package messagesmatcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
)

type MessagesMatcher struct {
	regexpMatcher *regexp.Regexp
}

func (m MessagesMatcher) MatchMessage(message message.Message) bool {
	isMatch := m.regexpMatcher.MatchString(strings.ToLower(message.MessageText))
	return isMatch
}

func NewMessagesMatcher(templateMatcher string) (MessagesMatcher, error) {
	regexpMatcher, err := regexp.Compile(templateMatcher)
	if err != nil {
		return MessagesMatcher{}, fmt.Errorf("error when compile regexp: %q", err)
	}
	return MessagesMatcher{regexpMatcher: regexpMatcher}, nil
}
