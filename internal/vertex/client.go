package vertex

import (
	"context"
	"log"

	"cloud.google.com/go/vertexai/genai"
)

func Client(projectID, region string, ctx context.Context) *genai.Client {
	client, err := genai.NewClient(ctx, projectID, region)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
