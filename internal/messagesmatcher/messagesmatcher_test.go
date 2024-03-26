package messagesmatcher

import (
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/stretchr/testify/assert"
)

func makeMessage(text string) message.Message {
	return message.Message{
		MessageText: text,
	}
}

func TestMessagesMatcher(t *testing.T) {
	matcher, err := NewMessagesMatcher(".*jira!.*")
	assert.Nil(t, err)
	assert.True(t, matcher.MatchMessage(makeMessage("hello, jira!")))
	assert.False(t, matcher.MatchMessage(makeMessage("hello, jira")))
	assert.True(t, matcher.MatchMessage(makeMessage("some text jira! other text")))
	assert.True(t, matcher.MatchMessage(makeMessage("some text jira!other text")))
	assert.True(t, matcher.MatchMessage(makeMessage("jIrA!")))
	assert.True(t, matcher.MatchMessage(makeMessage("JIRA!")))
	assert.True(t, matcher.MatchMessage(makeMessage("Jira! end of message")))
	assert.False(t, matcher.MatchMessage(makeMessage("random text")))
	assert.False(t, matcher.MatchMessage(makeMessage("there is a bug here. Can we start a task in jira?")))
	assert.True(t, matcher.MatchMessage(makeMessage("there is a bug here. Can we start a task in jira!?")))
	assert.False(t, matcher.MatchMessage(makeMessage("j i r A !")))
}

func TestBadReqexp(t *testing.T) {
	_, err := NewMessagesMatcher("+++")
	assert.NotNil(t, err)
	_, err = NewMessagesMatcher("**")
	assert.NotNil(t, err)
	_, err = NewMessagesMatcher("*aba")
	assert.NotNil(t, err)
}
