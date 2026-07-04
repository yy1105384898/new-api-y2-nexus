package setting

import (
	"strconv"
	"strings"
)

var CheckSensitiveEnabled = true
var CheckSensitiveOnPromptEnabled = true

// LocalSensitivePromptBlockEnabled 本地敏感词前置拦截：关闭后不检查词表、直接转发上游。
var LocalSensitivePromptBlockEnabled = true

// SensitiveReviewWhitelistUserIds 审查白名单用户：不受 LocalSensitivePromptBlockEnabled 控制，本地审查失败仍扣费。
var SensitiveReviewWhitelistUserIds = map[int]struct{}{}

//var CheckSensitiveOnCompletionEnabled = true

// StopOnSensitiveEnabled 如果检测到敏感词，是否立刻停止生成，否则替换敏感词
var StopOnSensitiveEnabled = true

// StreamCacheQueueLength 流模式缓存队列长度，0表示无缓存
var StreamCacheQueueLength = 0

// SensitiveWords 敏感词
// var SensitiveWords []string
var SensitiveWords = []string{
	"test_sensitive",
}

func SensitiveWordsToString() string {
	return strings.Join(SensitiveWords, "\n")
}

func SensitiveWordsFromString(s string) {
	SensitiveWords = []string{}
	sw := strings.Split(s, "\n")
	for _, w := range sw {
		w = strings.TrimSpace(w)
		if w != "" {
			SensitiveWords = append(SensitiveWords, w)
		}
	}
}

func SensitiveReviewWhitelistToString() string {
	if len(SensitiveReviewWhitelistUserIds) == 0 {
		return ""
	}
	ids := make([]string, 0, len(SensitiveReviewWhitelistUserIds))
	for id := range SensitiveReviewWhitelistUserIds {
		ids = append(ids, strconv.Itoa(id))
	}
	return strings.Join(ids, "\n")
}

func SensitiveReviewWhitelistFromString(s string) {
	SensitiveReviewWhitelistUserIds = map[int]struct{}{}
	for _, part := range strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n' || r == ',' || r == ' ' || r == '\t'
	}) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.Atoi(part)
		if err != nil || id <= 0 {
			continue
		}
		SensitiveReviewWhitelistUserIds[id] = struct{}{}
	}
}

func IsSensitiveReviewWhitelistUser(userId int) bool {
	if userId <= 0 {
		return false
	}
	_, ok := SensitiveReviewWhitelistUserIds[userId]
	return ok
}

// ShouldChargeOnLocalSensitiveRejection 白名单用户本地审查失败仍扣费（须先预扣再拦截）。
func ShouldChargeOnLocalSensitiveRejection(userId int) bool {
	return IsSensitiveReviewWhitelistUser(userId)
}

func promptSensitiveBaseEnabled() bool {
	return CheckSensitiveEnabled && CheckSensitiveOnPromptEnabled
}

func ShouldCheckPromptSensitive() bool {
	return ShouldCheckPromptSensitiveForUser(0)
}

// ShouldCheckPromptSensitiveForUser 白名单用户不受 LocalSensitivePromptBlockEnabled 全局开关影响。
func ShouldCheckPromptSensitiveForUser(userId int) bool {
	if !promptSensitiveBaseEnabled() {
		return false
	}
	if IsSensitiveReviewWhitelistUser(userId) {
		return true
	}
	return LocalSensitivePromptBlockEnabled
}

//func ShouldCheckCompletionSensitive() bool {
//	return CheckSensitiveEnabled && CheckSensitiveOnCompletionEnabled
//}
