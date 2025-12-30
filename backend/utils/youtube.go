package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// YouTubeURLPatterns contains regex patterns to extract YouTube video ID
var YouTubeURLPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`youtube\.com/shorts/([a-zA-Z0-9_-]{11})`),
}

// ExtractYouTubeID extracts the YouTube video ID from various URL formats
func ExtractYouTubeID(videoURL string) (string, error) {
	for _, pattern := range YouTubeURLPatterns {
		matches := pattern.FindStringSubmatch(videoURL)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("invalid YouTube URL: %s", videoURL)
}

// ValidateYouTubeURL validates if the URL is a valid YouTube URL
func ValidateYouTubeURL(videoURL string) bool {
	_, err := ExtractYouTubeID(videoURL)
	return err == nil
}

// GetYouTubeThumbnailURL returns the thumbnail URL for a YouTube video
func GetYouTubeThumbnailURL(videoID string, quality string) string {
	switch quality {
	case "high":
		return fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID)
	case "medium":
		return fmt.Sprintf("https://img.youtube.com/vi/%s/mqdefault.jpg", videoID)
	case "low":
		return fmt.Sprintf("https://img.youtube.com/vi/%s/sddefault.jpg", videoID)
	case "max":
		return fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
	default:
		return fmt.Sprintf("https://img.youtube.com/vi/%s/mqdefault.jpg", videoID)
	}
}

// YouTubeMetadata represents metadata fetched from noembed API
type YouTubeMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
	Thumbnail   string `json:"thumbnail_url"`
	Duration    string `json:"duration,omitempty"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	HTML        string `json:"html"`
}

// FetchYouTubeMetadata fetches video metadata from noembed API
func FetchYouTubeMetadata(videoURL string) (*YouTubeMetadata, error) {
	videoID, err := ExtractYouTubeID(videoURL)
	if err != nil {
		return nil, err
	}

	// Use noembed API to fetch metadata
	apiURL := fmt.Sprintf("https://noembed.com/embed?url=https://www.youtube.com/watch?v=%s", videoID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var metadata YouTubeMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// If noembed doesn't return description, fetch it differently
	if metadata.Description == "" {
		metadata.Description = getYouTubeDescription(videoID)
	}

	// Use YouTube's thumbnail if noembed doesn't provide one
	if metadata.Thumbnail == "" {
		metadata.Thumbnail = GetYouTubeThumbnailURL(videoID, "high")
	}

	return &metadata, nil
}

// FormatDuration formats seconds to human-readable duration
func FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

// ParseYouTubeDuration parses YouTube duration format (PT#H#M#S) to seconds
func ParseYouTubeDuration(duration string) int {
	duration = strings.TrimPrefix(duration, "PT")

	hours := 0
	minutes := 0
	seconds := 0

	if idx := strings.Index(duration, "H"); idx != -1 {
		fmt.Sscanf(duration[:idx], "%d", &hours)
		duration = duration[idx+1:]
	}

	if idx := strings.Index(duration, "M"); idx != -1 {
		fmt.Sscanf(duration[:idx], "%d", &minutes)
		duration = duration[idx+1:]
	}

	if idx := strings.Index(duration, "S"); idx != -1 {
		fmt.Sscanf(duration[:idx], "%d", &seconds)
	}

	return hours*3600 + minutes*60 + seconds
}

// getYouTubeDescription attempts to get video description (placeholder)
// YouTube doesn't provide a free API to get descriptions without API key
func getYouTubeDescription(videoID string) string {
	// This would require YouTube Data API with API key
	// For now, return empty string
	return ""
}

// BuildYouTubeURL constructs a YouTube URL from video ID
func BuildYouTubeURL(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
}

// GetYouTubeEmbedURL returns the embed URL for a YouTube video
func GetYouTubeEmbedURL(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/embed/%s", videoID)
}

// GetYouTubeShortsURL returns the Shorts URL for a YouTube video
func GetYouTubeShortsURL(videoID string) string {
	return fmt.Sprintf("https://www.youtube.com/shorts/%s", videoID)
}

// ExtractVideoIDFromURL is a more robust version of ExtractYouTubeID
func ExtractVideoIDFromURL(inputURL string) (string, error) {
	// First try to parse as a URL
	u, err := url.Parse(inputURL)
	if err == nil {
		// If it's a youtu.be short URL
		if u.Host == "youtu.be" {
			path := strings.TrimPrefix(u.Path, "/")
			if len(path) == 11 {
				return path, nil
			}
		}

		// If it's youtube.com
		if u.Host == "www.youtube.com" || u.Host == "youtube.com" {
			queryParams := u.Query()

			// Check for v parameter
			if v := queryParams.Get("v"); v != "" {
				return v, nil
			}

			// Check for embed path
			if strings.HasPrefix(u.Path, "/embed/") {
				id := strings.TrimPrefix(u.Path, "/embed/")
				return id, nil
			}

			// Check for shorts path
			if strings.HasPrefix(u.Path, "/shorts/") {
				id := strings.TrimPrefix(u.Path, "/shorts/")
				return id, nil
			}
		}
	}

	// Fall back to regex matching
	return ExtractYouTubeID(inputURL)
}
