package npms_test

import (
	"github.com/rjz/npms"
	"testing"
)

func expectQueryString(t *testing.T, desc, actual, expected string) {
	if actual != expected {
		t.Errorf("%s: Got '%s' (expected '%s')", desc, actual, expected)
	}
}

func TestQuery(t *testing.T) {
	expectQueryString(t, "q only", npms.SearchQuery("hello-world", nil), "hello-world")

	expectQueryString(t, "q + quals", npms.SearchQuery("hello-world", &npms.SearchQualifiers{
		Author:     "rjz",
		Maintainer: "rjz",
	}), "hello-world author:rjz,maintainer:rjz")

	// Author            string
	// Maintainer        string
	// Keywords          []string
	// BoostExact        bool
	// ScoreEffect       *float32
	// QualityWeight     *float32
	// PopularityWeight  *float32
	// MaintenanceWeight *float32
	expectQueryString(t, "quals only", npms.SearchQuery("", &npms.SearchQualifiers{
		Author:     "rjz",
		Maintainer: "rjz",
	}), "author:rjz,maintainer:rjz")

	expectQueryString(t, "kitchen sink", npms.SearchQuery("fzbz", &npms.SearchQualifiers{
		Author:            "rjz",
		BoostExact:        true,
		Keywords:          []string{"fz", "-bz", "fizz buzz"},
		Maintainer:        "rjz",
		MaintenanceWeight: npms.QueryFloat(0.50),
		PopularityWeight:  npms.QueryFloat(0.61),
		QualityWeight:     npms.QueryFloat(98.72),
		ScoreEffect:       npms.QueryFloat(0.83),
	}), "fzbz author:rjz,boost-exact:true,keywords:fz,-bz,fizz buzz,maintainer:rjz,maintenance-weight:0.50,popularity-weight:0.61,quality-weight:98.72,score-effect:0.83")

	expectQueryString(t, "no BoostExact", npms.SearchQuery("fzbz", &npms.SearchQualifiers{
		BoostExact: false,
	}), "fzbz")

	expectQueryString(t, "filters:not", npms.SearchQuery("fzbz", &npms.SearchQualifiers{
		Filters: npms.NotDeprecated | npms.NotUnstable | npms.NotInsecure,
	}), "fzbz not:deprecated,not:insecure,not:unstable")

	expectQueryString(t, "filters:is", npms.SearchQuery("fzbz", &npms.SearchQualifiers{
		Filters: npms.IsDeprecated | npms.IsUnstable | npms.IsInsecure,
	}), "fzbz is:deprecated,is:insecure,is:unstable")
}
