package model

// Stub constants and functions for removed LLM/billing features
// These allow auth-related code to compile while those features are removed

const (
	LogTypeSystem = 0
)

// RecordLog is a stub - logging functionality removed
func RecordLog(userId int, logType int, message string) {
	// No-op stub
}

// increaseTokenQuota is a stub - token quota removed
func increaseTokenQuota(tokenId int, quota int) error {
	return nil
}

// updateChannelUsedQuota is a stub - channel quota removed
func updateChannelUsedQuota(channelId int, quota int) error {
	return nil
}
