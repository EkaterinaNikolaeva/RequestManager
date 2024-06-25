package messagesmatcher

import (
	"context"
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
	tests := map[string]struct {
		msg      string
		expected bool
	}{
		"simple": {
			msg:      "hello, jira!",
			expected: true,
		},
		"without !": {
			msg:      "hello, jira",
			expected: false,
		},
		"in the middle": {
			msg:      "some text jira! other text!",
			expected: true,
		},
		"in the middle without space": {
			msg:      "some text jira!other text!",
			expected: true,
		},
		"upper and lower case": {
			msg:      "jIrA!",
			expected: true,
		},
		"all upper": {
			msg:      "jIrA!",
			expected: true,
		},
		"without pattern": {
			msg:      "random text",
			expected: false,
		},
		"real case false": {
			msg:      "there is a bug here. Can we start a task in jira?",
			expected: false,
		},
		"real case true": {
			msg:      "there is a bug here. Can we start a task in jira!?",
			expected: true,
		},
		"with spaces": {
			msg:      "j i r A !",
			expected: false,
		},
	}
	matcher, err := NewMessagesMatcher(".*jira!.*")
	if err != nil {
		t.Errorf("TestMessagesMatcher error %q", err)
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, matcher.MatchMessage(context.Background(), makeMessage(tc.msg)))
		})
	}
}

func TestBadReqexp(t *testing.T) {
	tests := map[string]struct {
		regexp string
	}{
		"all +": {
			regexp: "+++",
		},
		"all *": {
			regexp: "**",
		},
		"starts with *": {
			regexp: "*aba",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewMessagesMatcher(tc.regexp)
			assert.NotNil(t, err)
		})
	}
}
