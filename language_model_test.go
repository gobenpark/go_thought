package main

import (
	"context"
	"html/template"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTemplate(t *testing.T) {
	tm := template.New("template")
	temp, err := tm.Parse(`Hello {{.Name}}
{{if hasKey .Ben }}
hello
{{- end}}

`)
	require.NoError(t, err)
	temp.Execute(os.Stdout, map[string]string{
		"Name": "ben",
		"Ben":  "",
	})
}

func TestOpenAIProvider_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	op := NewMockProvider(ctrl)

	expectOpenAIMessage := []Message{
		{
			Role:    "system",
			Message: "you are a ai",
		},
		{
			Role:    "user",
			Message: "what time is today?",
		},
	}

	op.EXPECT().Generate(gomock.Any(), nil, expectOpenAIMessage).Return([]ResponseMessage{
		{
			Message: "time is 00:00",
		},
	}, nil)
	model := NewLanguageModel(op)

	res, err := model.SystemPrompt("you are a ai").HumanPrompt("what time is today?").Q(context.TODO())
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, "time is 00:00", res.Message)
}
