package vertex

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

func Send(ctx context.Context, cs *genai.ChatSession, msg string) *genai.GenerateContentResponse {
	res, err := cs.SendMessage(ctx, genai.Text(msg))
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func Response(resp *genai.GenerateContentResponse) string {
	var s []string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			s = append(s, fmt.Sprintf("%#v", part))
		}
	}
	return strings.Join(s[:], ",")

}
