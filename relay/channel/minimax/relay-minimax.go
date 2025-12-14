package minimax

import (
	"fmt"

	channelconstant "github.com/karlchow/auth-gateway/constant"
	relaycommon "github.com/karlchow/auth-gateway/relay/common"
	"github.com/karlchow/auth-gateway/relay/constant"
)

func GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	baseUrl := info.ChannelBaseUrl
	if baseUrl == "" {
		baseUrl = channelconstant.ChannelBaseURLs[channelconstant.ChannelTypeMiniMax]
	}

	switch info.RelayMode {
	case constant.RelayModeChatCompletions:
		return fmt.Sprintf("%s/v1/text/chatcompletion_v2", baseUrl), nil
	case constant.RelayModeAudioSpeech:
		return fmt.Sprintf("%s/v1/t2a_v2", baseUrl), nil
	default:
		return "", fmt.Errorf("unsupported relay mode: %d", info.RelayMode)
	}
}
