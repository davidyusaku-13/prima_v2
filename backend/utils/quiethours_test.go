package utils

import (
	"testing"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
)

// Helper function to create pointer to int
func intPtr(i int) *int {
	return &i
}

func TestIsQuietHours_At2100_ReturnsTrue(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 21:00 WIB = 14:00 UTC
	testTime := time.Date(2025, 12, 29, 14, 0, 0, 0, time.UTC)

	if !IsQuietHours(testTime, cfg) {
		t.Error("Expected 21:00 WIB to be in quiet hours")
	}
}

func TestIsQuietHours_At0559_ReturnsTrue(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 05:59 WIB = 22:59 UTC (previous day)
	testTime := time.Date(2025, 12, 28, 22, 59, 0, 0, time.UTC)

	if !IsQuietHours(testTime, cfg) {
		t.Error("Expected 05:59 WIB to be in quiet hours")
	}
}

func TestIsQuietHours_At0600_ReturnsFalse(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 06:00 WIB = 23:00 UTC (previous day)
	testTime := time.Date(2025, 12, 28, 23, 0, 0, 0, time.UTC)

	if IsQuietHours(testTime, cfg) {
		t.Error("Expected 06:00 WIB to NOT be in quiet hours")
	}
}

func TestIsQuietHours_At2059_ReturnsFalse(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 20:59 WIB = 13:59 UTC
	testTime := time.Date(2025, 12, 29, 13, 59, 0, 0, time.UTC)

	if IsQuietHours(testTime, cfg) {
		t.Error("Expected 20:59 WIB to NOT be in quiet hours")
	}
}

func TestIsQuietHours_AtMidnight_ReturnsTrue(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 00:00 WIB = 17:00 UTC (previous day)
	testTime := time.Date(2025, 12, 28, 17, 0, 0, 0, time.UTC)

	if !IsQuietHours(testTime, cfg) {
		t.Error("Expected midnight WIB to be in quiet hours")
	}
}

func TestIsQuietHours_At1200_ReturnsFalse(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 12:00 WIB = 05:00 UTC
	testTime := time.Date(2025, 12, 29, 5, 0, 0, 0, time.UTC)

	if IsQuietHours(testTime, cfg) {
		t.Error("Expected 12:00 WIB to NOT be in quiet hours")
	}
}

func TestGetNextActiveTime_DuringQuietHours_BeforeMidnight_ReturnsToday0600(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 22:00 WIB on Dec 29 = 15:00 UTC on Dec 29
	testTime := time.Date(2025, 12, 29, 15, 0, 0, 0, time.UTC)

	nextActive := GetNextActiveTime(testTime, cfg)

	// Expected: 06:00 WIB on Dec 30 = 23:00 UTC on Dec 29
	expected := time.Date(2025, 12, 29, 23, 0, 0, 0, time.UTC)

	if !nextActive.Equal(expected) {
		t.Errorf("Expected next active time %v, got %v", expected, nextActive)
	}
}

func TestGetNextActiveTime_DuringQuietHours_AfterMidnight_ReturnsToday0600(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 03:00 WIB on Dec 30 = 20:00 UTC on Dec 29
	testTime := time.Date(2025, 12, 29, 20, 0, 0, 0, time.UTC)

	nextActive := GetNextActiveTime(testTime, cfg)

	// Expected: 06:00 WIB on Dec 30 = 23:00 UTC on Dec 29
	expected := time.Date(2025, 12, 29, 23, 0, 0, 0, time.UTC)

	if !nextActive.Equal(expected) {
		t.Errorf("Expected next active time %v, got %v", expected, nextActive)
	}
}

func TestGetNextActiveTime_After0600_ReturnsTomorrow0600(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 10:00 WIB on Dec 29 = 03:00 UTC on Dec 29
	testTime := time.Date(2025, 12, 29, 3, 0, 0, 0, time.UTC)

	nextActive := GetNextActiveTime(testTime, cfg)

	// Expected: 06:00 WIB on Dec 30 = 23:00 UTC on Dec 29
	expected := time.Date(2025, 12, 29, 23, 0, 0, 0, time.UTC)

	if !nextActive.Equal(expected) {
		t.Errorf("Expected next active time %v, got %v", expected, nextActive)
	}
}

func TestGetNextActiveTime_At0600Exactly_ReturnsTomorrow0600(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 06:00 WIB on Dec 29 = 23:00 UTC on Dec 28
	testTime := time.Date(2025, 12, 28, 23, 0, 0, 0, time.UTC)

	nextActive := GetNextActiveTime(testTime, cfg)

	// Expected: 06:00 WIB on Dec 30 = 23:00 UTC on Dec 29
	expected := time.Date(2025, 12, 29, 23, 0, 0, 0, time.UTC)

	if !nextActive.Equal(expected) {
		t.Errorf("Expected next active time %v, got %v", expected, nextActive)
	}
}

func TestGetTimezoneLocation_WIB(t *testing.T) {
	loc := GetTimezoneLocation("WIB")
	if loc.String() != "WIB" {
		t.Errorf("Expected WIB location, got %s", loc.String())
	}

	// Verify offset is +7 hours
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	localTime := testTime.In(loc)
	if localTime.Hour() != 7 {
		t.Errorf("Expected hour 7 in WIB, got %d", localTime.Hour())
	}
}

func TestGetTimezoneLocation_WITA(t *testing.T) {
	loc := GetTimezoneLocation("WITA")
	if loc.String() != "WITA" {
		t.Errorf("Expected WITA location, got %s", loc.String())
	}

	// Verify offset is +8 hours
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	localTime := testTime.In(loc)
	if localTime.Hour() != 8 {
		t.Errorf("Expected hour 8 in WITA, got %d", localTime.Hour())
	}
}

func TestGetTimezoneLocation_WIT(t *testing.T) {
	loc := GetTimezoneLocation("WIT")
	if loc.String() != "WIT" {
		t.Errorf("Expected WIT location, got %s", loc.String())
	}

	// Verify offset is +9 hours
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	localTime := testTime.In(loc)
	if localTime.Hour() != 9 {
		t.Errorf("Expected hour 9 in WIT, got %d", localTime.Hour())
	}
}

func TestGetTimezoneLocation_Unknown_DefaultsToWIB(t *testing.T) {
	loc := GetTimezoneLocation("UNKNOWN")
	if loc.String() != "WIB" {
		t.Errorf("Expected default WIB location for unknown timezone, got %s", loc.String())
	}
}

func TestIsQuietHours_WithWITA(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WITA",
	}

	// 21:00 WITA = 13:00 UTC
	testTime := time.Date(2025, 12, 29, 13, 0, 0, 0, time.UTC)

	if !IsQuietHours(testTime, cfg) {
		t.Error("Expected 21:00 WITA to be in quiet hours")
	}
}

func TestIsQuietHours_WithWIT(t *testing.T) {
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(21),
		EndHour:   intPtr(6),
		Timezone:  "WIT",
	}

	// 21:00 WIT = 12:00 UTC
	testTime := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)

	if !IsQuietHours(testTime, cfg) {
		t.Error("Expected 21:00 WIT to be in quiet hours")
	}
}

func TestIsQuietHours_SameStartAndEnd_ReturnsFalse(t *testing.T) {
	// When start and end are the same, quiet hours is effectively disabled
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(6),
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// Any time should return false when start == end
	testTime := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)

	if IsQuietHours(testTime, cfg) {
		t.Error("Expected same start/end quiet hours to return false")
	}

	// Test at midnight too
	testTime = time.Date(2025, 12, 29, 0, 0, 0, 0, time.UTC)
	if IsQuietHours(testTime, cfg) {
		t.Error("Expected same start/end quiet hours to return false at midnight")
	}
}

func TestIsQuietHours_MidnightStartHour_Works(t *testing.T) {
	// Test that StartHour=0 (midnight) works correctly and is not treated as "unset"
	cfg := &config.QuietHoursConfig{
		StartHour: intPtr(0),  // Midnight
		EndHour:   intPtr(6),
		Timezone:  "WIB",
	}

	// 00:30 WIB = 17:30 UTC (previous day)
	testTime := time.Date(2025, 12, 28, 17, 30, 0, 0, time.UTC)

	if !IsQuietHours(testTime, cfg) {
		t.Error("Expected 00:30 WIB to be in quiet hours when StartHour=0")
	}

	// 23:00 WIB = 16:00 UTC - should NOT be in quiet hours
	testTime = time.Date(2025, 12, 28, 16, 0, 0, 0, time.UTC)

	if IsQuietHours(testTime, cfg) {
		t.Error("Expected 23:00 WIB to NOT be in quiet hours when StartHour=0")
	}
}
