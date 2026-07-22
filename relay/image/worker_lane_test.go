package image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdobeDirectChannelIDsAreConfigurable(t *testing.T) {
	t.Setenv("IMAGE_ASYNC_ADOBE_CHANNEL_IDS", "75, 81;75 invalid")
	require.Equal(t, []int{75, 81}, AdobeDirectChannelIDs())
	require.True(t, IsAdobeDirectChannel(75))
	require.True(t, IsAdobeDirectChannel(81))
	require.False(t, IsAdobeDirectChannel(77))
}

func TestImageDispatcherQueueKeysWorkBeforeWorkerStartup(t *testing.T) {
	notify, dedup := imageDispatcherQueueKeys(&imageDispatcher)
	require.Equal(t, imageTaskNotifyQueue, notify)
	require.Equal(t, imageTaskNotifyDedup, dedup)

	notify, dedup = imageDispatcherQueueKeys(&adobeImageDispatcher)
	require.Equal(t, adobeTaskNotifyQueue, notify)
	require.Equal(t, adobeTaskNotifyDedup, dedup)
}
