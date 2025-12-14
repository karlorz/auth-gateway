package model

// Stub constants and functions for removed LLM/billing features
// These allow auth-related code to compile while those features are removed

const (
	LogTypeSystem = 0
	LogTypeManage = 1
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

// Token is a stub struct for removed API token functionality
type Token struct {
	Id     int
	UserId int
	Key    string
}

// GetGroupEnabledModels is a stub - models functionality removed
func GetGroupEnabledModels(group string) []string {
	return []string{}
}

// Redeem is a stub - redemption functionality removed
func Redeem(code string, userId int) (int, error) {
	return 0, nil
}
