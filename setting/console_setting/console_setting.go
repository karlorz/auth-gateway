package console_setting

// ConsoleSetting holds dashboard/console configuration
type ConsoleSetting struct {
	ApiInfoEnabled      bool `json:"api_info_enabled"`
	UptimeKumaEnabled   bool `json:"uptime_kuma_enabled"`
	AnnouncementsEnabled bool `json:"announcements_enabled"`
	FAQEnabled          bool `json:"faq_enabled"`
	HeaderNavModules    string `json:"HeaderNavModules"`
	SidebarModulesAdmin string `json:"SidebarModulesAdmin"`
	ApiInfo             []interface{} `json:"api_info"`
	Announcements       []interface{} `json:"announcements"`
	FAQ                 []interface{} `json:"faq"`
}

var defaultConsoleSetting = ConsoleSetting{
	ApiInfoEnabled:      true,
	UptimeKumaEnabled:   true,
	AnnouncementsEnabled: true,
	FAQEnabled:          true,
}

func GetConsoleSetting() *ConsoleSetting {
	return &defaultConsoleSetting
}

func SetConsoleSetting(cs *ConsoleSetting) {
	defaultConsoleSetting = *cs
}
