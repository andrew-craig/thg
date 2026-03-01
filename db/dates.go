package db

import (
	"fmt"
	"time"
)

// DecodeDate converts a Things bit-packed date integer to a time.Time.
// Encoding: (year << 16) | (month << 12) | (day << 7)
func DecodeDate(v int64) time.Time {
	year := int(v >> 16)
	month := time.Month((v >> 12) & 0xF)
	day := int((v >> 7) & 0x1F)
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// EncodeDate converts a time.Time to a Things bit-packed date integer.
func EncodeDate(t time.Time) int64 {
	return int64(t.Year())<<16 | int64(t.Month())<<12 | int64(t.Day())<<7
}

// FormatDate returns a human-readable date string, or empty if the value is 0.
func FormatDate(v int64) string {
	if v == 0 {
		return ""
	}
	t := DecodeDate(v)
	return t.Format("2006-01-02")
}

// TodayEncoded returns today's date as a Things bit-packed integer.
func TodayEncoded() int64 {
	return EncodeDate(time.Now())
}

// FormatTimestamp formats a Unix timestamp (float seconds) as a date string.
func FormatTimestamp(ts float64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04")
}

// ParseDateArg parses a user-provided date string (YYYY-MM-DD) into a Things date int.
func ParseDateArg(s string) (int64, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, fmt.Errorf("invalid date %q (use YYYY-MM-DD)", s)
	}
	return EncodeDate(t), nil
}
