package main

import (
	"reflect"
	"sort"
	"testing"

	"github.com/posteo/go-agentx/value"
)

func TestOIDSort(t *testing.T) {
	in := []value.OID{
		value.MustParseOID("1.3.6"),
		value.MustParseOID("1.3.6.3"),
		value.MustParseOID("1.3.6.3.1"),
		value.MustParseOID("1.3.6.3"),
		value.MustParseOID("1.3.6"),
		value.MustParseOID("1.3.6.1.3"),
		value.MustParseOID("1.3.6.10.2"),
		value.MustParseOID("1.3.6.100.1"),
		value.MustParseOID("1.3.6.10.1"),
		value.MustParseOID("1.3.6.1.2"),
		value.MustParseOID("1.3.6.2.2"),
		value.MustParseOID("1.3.6.1.3"),
		value.MustParseOID("1.3.6.2.1"),
	}

	out := []value.OID{
		value.MustParseOID("1.3.6"),
		value.MustParseOID("1.3.6"),
		value.MustParseOID("1.3.6.1.2"),
		value.MustParseOID("1.3.6.1.3"),
		value.MustParseOID("1.3.6.1.3"),
		value.MustParseOID("1.3.6.2.1"),
		value.MustParseOID("1.3.6.2.2"),
		value.MustParseOID("1.3.6.3"),
		value.MustParseOID("1.3.6.3"),
		value.MustParseOID("1.3.6.3.1"),
		value.MustParseOID("1.3.6.10.1"),
		value.MustParseOID("1.3.6.10.2"),
		value.MustParseOID("1.3.6.100.1"),
	}

	sort.Sort(OIDSorter(in))
	if !reflect.DeepEqual(in, out) {
		t.Errorf("Error sorting oids. Expected %+v. Got %+v", out, in)
	}
}

func TestOIDComparison(t *testing.T) {
	test_data := []struct {
		left      string
		right     string
		less_than bool
	}{
		{"1.3.6.1", "1.3.6.2", true},
		{"1.3.6.1", "1.3.6.10", true},
		{"1.3.6.1.1", "1.3.6.10.1", true},
		{"1.3.6.1.1", "1.3.6", false},
		{"1.30.6.1.1", "1.300.6", true},
	}

	for _, test_row := range test_data {
		if OIDLessThan(value.MustParseOID(test_row.left), value.MustParseOID(test_row.right)) != test_row.less_than {
			t.Errorf("%v < %v was expected to be %v", test_row.left, test_row.right, test_row.less_than)
		}

		if OIDLessThan(value.MustParseOID(test_row.right), value.MustParseOID(test_row.left)) == test_row.less_than {
			t.Errorf("%v < %v was expected to be %v", test_row.right, test_row.left, !test_row.less_than)
		}

		if OIDGreaterThan(value.MustParseOID(test_row.left), value.MustParseOID(test_row.right)) == test_row.less_than {
			t.Errorf("%v > %v was expected to be %v", test_row.left, test_row.right, !test_row.less_than)
		}

		if OIDGreaterThan(value.MustParseOID(test_row.right), value.MustParseOID(test_row.left)) != test_row.less_than {
			t.Errorf("%v > %v was expected to be %v", test_row.right, test_row.left, test_row.less_than)
		}
	}
}
