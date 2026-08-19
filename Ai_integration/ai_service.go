package aiintegration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Prkarz/course-tracker/models"
	"google.golang.org/genai"
)

func Client_Ai_Init() (*genai.Client, error) {
	ctx := context.Background()
	api_key := os.Getenv("GOOGLE_API_KEY")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting with AI: %w", err)
	}
	return client, nil
}

func CourseInsights(ctx context.Context, client *genai.Client, summary string) (*models.AIresponseCatcher, error) {
	prompt := fmt.Sprintf("Generate a Detailed Summary of the course: %s. You MUST return a JSON object containing exactly two keys. The first key must be named 'ai_summary' and the second key must be named 'course_tags'.", summary)

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return nil, err
	}
	var result models.AIresponseCatcher
	err = json.Unmarshal([]byte(resp.Text()), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing AI response: %w", err)
	}

	return &result, nil
}
