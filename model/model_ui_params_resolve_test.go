package model

import "testing"

func TestProfileToDocumentImageHints(t *testing.T) {
	doc, err := profileToDocument(ModelUiParamProfile{
		ProfileId:   "image-tpl-test",
		Capability:  ModelUiParamCapabilityImage,
		Params:      `{"quality":{"enabled":false}}`,
		Hints:       `[{"text":"示例 hint"}]`,
		OptionRules: "[]",
	})
	if err != nil {
		t.Fatalf("profileToDocument() error = %v", err)
	}
	hints, ok := doc["hints"].([]interface{})
	if !ok || len(hints) != 1 {
		t.Fatalf("hints = %#v, want one hint", doc["hints"])
	}
}
