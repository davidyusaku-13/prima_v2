package utils

import (
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
)

// WIBLocation is the Western Indonesia Time timezone (UTC+7)
var WIBLocation = time.FixedZone("WIB", 7*60*60)

// WITALocation is the Central Indonesia Time timezone (UTC+8)
var WITALocation = time.FixedZone("WITA", 8*60*60)

// WITLocation is the Eastern Indonesia Time timezone (UTC+9)
var WITLocation = time.FixedZone("WIT", 9*60*60)

// GetTimezoneLocation returns the time.Location for the given timezone string
func GetTimezoneLocation(tz string) *time.Location {
	switch tz {
	case "WITA":
		return WITALocation
	case "WIT":
		return WITLocation
	default:
		return WIBLocation // Default to WIB
	}
}

// TimezoneOffset returns the UTC offset in seconds for the given timezone
func TimezoneOffset(tz string) int {
	switch tz {
	case "WITA":
		return 8 * 60 * 60
	case "WIT":
		return 9 * 60 * 60
	default:
		return 7 * 60 * 60 // WIB
	}
}

// IsQuietHours checks if the given time falls within quiet hours
// Handles both cases:
// - Quiet hours spanning midnight (e.g., 21:00 - 06:00): hour >= start OR hour < end
// - Quiet hours not spanning midnight (e.g., 00:00 - 06:00): hour >= start AND hour < end
func IsQuietHours(t time.Time, cfg *config.QuietHoursConfig) bool {
	startHour := cfg.GetStartHour()
	endHour := cfg.GetEndHour()

	// If start and end are the same, quiet hours is effectively disabled
	if startHour == endHour {
		return false
	}

	loc := GetTimezoneLocation(cfg.Timezone)
	localTime := t.In(loc)
	hour := localTime.Hour()

	// Check if quiet hours span midnight
	if startHour > endHour {
		// Quiet hours span midnight (e.g., 21:00 - 06:00)
		// Check if hour >= start OR hour < end
		return hour >= startHour || hour < endHour
	}
	// Quiet hours don't span midnight (e.g., 00:00 - 06:00)
	// Check if hour >= start AND hour < end
	return hour >= startHour && hour < endHour
}

// GetNextActiveTime returns the next time when active hours begin (EndHour)
// If current time is before EndHour today, returns today at EndHour
// If current time is at or after EndHour, returns tomorrow at EndHour
func GetNextActiveTime(t time.Time, cfg *config.QuietHoursConfig) time.Time {
	loc := GetTimezoneLocation(cfg.Timezone)
	localTime := t.In(loc)
	endHour := cfg.GetEndHour()

	// Set to EndHour (06:00) today in local timezone
	nextActive := time.Date(
		localTime.Year(), localTime.Month(), localTime.Day(),
		endHour, 0, 0, 0, loc,
	)

	// If current time is at or past EndHour, move to tomorrow
	if localTime.Hour() >= endHour {
		nextActive = nextActive.AddDate(0, 0, 1)
	}

	// Return in UTC for storage
	return nextActive.UTC()
}
