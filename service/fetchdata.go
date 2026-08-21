package service

import (
	"context"
	"errors"
	"os"

	"github.com/Prkarz/course-tracker/models"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func FetchPlayListVideos(listID string) ([]models.VideoData, error) {
	api_key := os.Getenv("GOOGLE_API_KEY")
	if api_key == "" {
		return nil, errors.New("Error making Client")
	}
	ctx := context.Background()

	services, err := youtube.NewService(ctx, option.WithAPIKey(api_key))
	if err != nil {
		return nil, err
	}

	call := services.PlaylistItems.List([]string{"snippet"}).PlaylistId(listID).MaxResults(50)
	resposnse, err := call.Do()
	if err != nil {
		return nil, err
	}
	var temp []models.VideoData
	for _, vid := range resposnse.Items {
		if vid.Snippet.Title == "Private video" || vid.Snippet.Title == "Deleted video" {
			continue
		}
		item := models.VideoData{
			VideoId:    vid.Snippet.ResourceId.VideoId,
			VideoTitle: vid.Snippet.Title,
			Index:      int(vid.Snippet.Position) + 1,
			Duration:   "TBD",
		}
		temp = append(temp, item)
	}
	return temp, nil
}
