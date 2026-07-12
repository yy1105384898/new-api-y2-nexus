package registry

import "testing"

func TestResolve(t *testing.T) {
	cases := []struct {
		origin   string
		upstream string
		want     Vendor
	}{
		{"manju-openai-sora2", "sora2", VendorManju},
		{"cy-sd1-seedance-2.0-fast-720p", "Seedance-2.0-720p", VendorSeedance},
		{"cy-sd4-seedance-2.0", "seedance-2.0", VendorSeedance},
		{"cy-sd2-seedance-2.0", "manxue-2.0", VendorSeedance},
		{"tengd-seedance-2.0", "manxue-2.0", VendorSeedance},
		{"cy-vid2-sora-2", "cy-vid2-sora-2", VendorChat},
		{"cy-sd1-grok-video", "grok-imagine-video", VendorChat},
		{"cy-gv1-grok-video", "grok-image-video", VendorGrok},
		{"cy-gv1-grok-video-1.5", "grok-video-1.5", VendorGrok},
		{"sora-2", "sora-2", VendorSora},
		{"grok-video", "grok-video", VendorSora},
	}
	for _, tc := range cases {
		if got := Resolve(tc.origin, tc.upstream); got != tc.want {
			t.Fatalf("Resolve(%q,%q)=%q want %q", tc.origin, tc.upstream, got, tc.want)
		}
	}
}

func TestResolve_CySd1NotTengda(t *testing.T) {
	if Resolve("cy-sd1-seedance-2.0-720p", "manxue-2.0") != VendorSeedance {
		t.Fatal("cy-sd1 should resolve to seedance even with tengda upstream name")
	}
}
