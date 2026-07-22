package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSameImagePublicOrigin(t *testing.T) {
	require.True(t, sameImagePublicOrigin(
		"https://tmp.cangyuansuanli.cn/gen-images/1/task/0.png",
		"https://tmp.cangyuansuanli.cn",
	))
	require.True(t, sameImagePublicOrigin(
		"https://cdn.example.com/media/gen-images/1/task/0.png",
		"https://cdn.example.com/media",
	))
	require.False(t, sameImagePublicOrigin(
		"https://tmp.cangyuansuanli.cn.evil.example/gen-images/1/task/0.png",
		"https://tmp.cangyuansuanli.cn",
	))
	require.False(t, sameImagePublicOrigin(
		"http://tmp.cangyuansuanli.cn/gen-images/1/task/0.png",
		"https://tmp.cangyuansuanli.cn",
	))
	require.False(t, sameImagePublicOrigin(
		"https://cdn.example.com/other/0.png",
		"https://cdn.example.com/media",
	))
}
