package service

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Prkarz/course-tracker/models"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var isoDurationRegex = regexp.MustCompile(`P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?`)

// parseISO8601Duration converts ISO 8601 (e.g. "PT1H15M33S", "PT15M33S", "PT45S", "P1DT2H") to "1:15:33" or "15:33"
func parseISO8601Duration(isoDuration string) string {
	isoDuration = strings.TrimSpace(isoDuration)
	if isoDuration == "" || isoDuration == "P0D" || isoDuration == "PT0S" {
		return "0:00"
	}

	matches := isoDurationRegex.FindStringSubmatch(isoDuration)
	if len(matches) == 0 {
		return "10:00"
	}

	var days, hours, minutes, seconds int
	if len(matches) > 1 && matches[1] != "" {
		fmt.Sscanf(matches[1], "%d", &days)
	}
	if len(matches) > 2 && matches[2] != "" {
		fmt.Sscanf(matches[2], "%d", &hours)
	}
	if len(matches) > 3 && matches[3] != "" {
		fmt.Sscanf(matches[3], "%d", &minutes)
	}
	if len(matches) > 4 && matches[4] != "" {
		fmt.Sscanf(matches[4], "%d", &seconds)
	}

	hours += days * 24

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// FetchVideoDurations queries the YouTube Videos API in batches of 50 to retrieve exact durations
func FetchVideoDurations(videoIDs []string) map[string]string {
	if len(videoIDs) == 0 {
		return make(map[string]string)
	}
	apiKey := strings.TrimSpace(os.Getenv("YOUTUBE_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	}
	if apiKey == "" {
		return make(map[string]string)
	}
	ctx := context.Background()
	services, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return make(map[string]string)
	}
	return fetchDurationsMap(services, videoIDs)
}

// fetchDurationsMap queries the YouTube Videos API in batches of 50 to retrieve real video durations
func fetchDurationsMap(services *youtube.Service, videoIDs []string) map[string]string {
	durMap := make(map[string]string)
	if len(videoIDs) == 0 {
		return durMap
	}

	for i := 0; i < len(videoIDs); i += 50 {
		end := i + 50
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		batch := videoIDs[i:end]
		call := services.Videos.List([]string{"contentDetails"}).Id(strings.Join(batch, ","))
		res, err := call.Do()
		if err != nil {
			continue
		}
		for _, item := range res.Items {
			if item.ContentDetails != nil {
				durMap[item.Id] = parseISO8601Duration(item.ContentDetails.Duration)
			}
		}
	}
	return durMap
}

func FetchPlayListVideos(listID string) ([]models.VideoData, error) {
	apiKey := strings.TrimSpace(os.Getenv("YOUTUBE_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	}
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY is not configured")
	}
	ctx := context.Background()

	services, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("create YouTube client: %w", err)
	}

	var rawItems []*youtube.PlaylistItem
	nextPageToken := ""

	for {
		call := services.PlaylistItems.List([]string{"snippet"}).PlaylistId(listID).MaxResults(50)
		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}
		response, err := call.Do()
		if err != nil {
			if len(rawItems) == 0 {
				return nil, fmt.Errorf("fetch playlist %q: %w", listID, err)
			}
			break
		}
		rawItems = append(rawItems, response.Items...)
		nextPageToken = response.NextPageToken
		if nextPageToken == "" || len(rawItems) >= 200 {
			break
		}
	}

	// 1. Collect and deduplicate unique video IDs
	seenVideoIDs := make(map[string]bool)
	var uniqueVideoIDs []string

	type basicVideo struct {
		videoID string
		title   string
	}
	var orderedList []basicVideo

	for _, vid := range rawItems {
		if vid.Snippet == nil || vid.Snippet.ResourceId == nil {
			continue
		}
		title := strings.TrimSpace(vid.Snippet.Title)
		if title == "" || title == "Private video" || title == "Deleted video" {
			continue
		}
		vID := strings.TrimSpace(vid.Snippet.ResourceId.VideoId)
		if vID == "" || seenVideoIDs[vID] {
			continue
		}
		seenVideoIDs[vID] = true
		uniqueVideoIDs = append(uniqueVideoIDs, vID)
		orderedList = append(orderedList, basicVideo{
			videoID: vID,
			title:   title,
		})
	}

	// 2. Fetch real durations from YouTube contentDetails API
	durationsMap := fetchDurationsMap(services, uniqueVideoIDs)

	// 3. Construct clean sequentially indexed VideoData list (#1, #2, #3...)
	var temp []models.VideoData
	for idx, vid := range orderedList {
		duration := durationsMap[vid.videoID]
		if duration == "" {
			duration = "10:00"
		}
		item := models.VideoData{
			VideoId:    vid.videoID,
			VideoTitle: vid.title,
			Index:      idx + 1, // Clean sequential 1-based index (#1, #2, #3...)
			Duration:   duration,
		}
		temp = append(temp, item)
	}

	return temp, nil
}
