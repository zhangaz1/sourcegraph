package licensing

import (
	"testing"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/license"
)

func TestCheckFeature(t *testing.T) {
	license := func(tags ...string) *Info { return &Info{Info: license.Info{Tags: tags}} }

	check := func(t *testing.T, feature Feature, info *Info, wantEnabled bool) {
		t.Helper()
		got := checkFeature(info, feature) == nil
		if got != wantEnabled {
			t.Errorf("got %v, want %v", got, wantEnabled)
		}
	}

	t.Run(string(FeatureACLs), func(t *testing.T) {
		check(t, FeatureACLs, nil, false)
		check(t, FeatureACLs, license(freeNoLicense.tag()), false)
		check(t, FeatureACLs, license("starter"), false)
		check(t, FeatureACLs, license(), true)
	})

	t.Run(string(FeatureExtensionRegistry), func(t *testing.T) {
		check(t, FeatureExtensionRegistry, nil, false)
		check(t, FeatureExtensionRegistry, license(freeNoLicense.tag()), false)
		check(t, FeatureExtensionRegistry, license("starter"), false)
		check(t, FeatureExtensionRegistry, license(), true)
	})

	t.Run(string(FeatureRemoteExtensionsAllowDisallow), func(t *testing.T) {
		check(t, FeatureRemoteExtensionsAllowDisallow, nil, false)
		check(t, FeatureRemoteExtensionsAllowDisallow, license(freeNoLicense.tag()), false)
		check(t, FeatureRemoteExtensionsAllowDisallow, license("starter"), false)
		check(t, FeatureRemoteExtensionsAllowDisallow, license(), true)
	})

	t.Run(string(FeatureMultipleExternalServices), func(t *testing.T) {
		t.Run("initialized pre-3.15", func(t *testing.T) {
			check(t, FeatureMultipleExternalServices, nil, true)
			check(t, FeatureMultipleExternalServices, license(freeNoLicense.tag()), true)
		})
		t.Run("initialized post-3.15", func(t *testing.T) {
			check(t, FeatureMultipleExternalServices, nil, false)
			check(t, FeatureMultipleExternalServices, license(freeNoLicense.tag()), false)
		})
		check(t, FeatureMultipleExternalServices, license("starter"), true)
		check(t, FeatureMultipleExternalServices, license(), true)
	})
}
