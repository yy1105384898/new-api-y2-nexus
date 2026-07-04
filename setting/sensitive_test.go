package setting

import "testing"

func TestShouldCheckPromptSensitiveForUser_WhitelistBypassesGlobalOff(t *testing.T) {
	prevGlobal := LocalSensitivePromptBlockEnabled
	prevEnabled := CheckSensitiveEnabled
	prevPrompt := CheckSensitiveOnPromptEnabled
	prevWhitelist := SensitiveReviewWhitelistUserIds
	t.Cleanup(func() {
		LocalSensitivePromptBlockEnabled = prevGlobal
		CheckSensitiveEnabled = prevEnabled
		CheckSensitiveOnPromptEnabled = prevPrompt
		SensitiveReviewWhitelistUserIds = prevWhitelist
	})

	LocalSensitivePromptBlockEnabled = false
	CheckSensitiveEnabled = true
	CheckSensitiveOnPromptEnabled = true
	SensitiveReviewWhitelistUserIds = map[int]struct{}{42: {}}

	if ShouldCheckPromptSensitiveForUser(1) {
		t.Fatal("expected non-whitelist user to skip check when global block disabled")
	}
	if !ShouldCheckPromptSensitiveForUser(42) {
		t.Fatal("expected whitelist user to keep local check when global block disabled")
	}
}

func TestSensitiveReviewWhitelistFromString(t *testing.T) {
	prev := SensitiveReviewWhitelistUserIds
	t.Cleanup(func() {
		SensitiveReviewWhitelistUserIds = prev
	})

	SensitiveReviewWhitelistFromString("1\n2, 3\t4")
	if !IsSensitiveReviewWhitelistUser(1) || !IsSensitiveReviewWhitelistUser(2) || !IsSensitiveReviewWhitelistUser(3) || !IsSensitiveReviewWhitelistUser(4) {
		t.Fatalf("unexpected whitelist: %#v", SensitiveReviewWhitelistUserIds)
	}
	if IsSensitiveReviewWhitelistUser(5) {
		t.Fatal("user 5 should not be whitelisted")
	}
}
