package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseLicenseStats(t *testing.T) {
	testCases := []struct {
		input string
		want  licenseStats
		ok    bool
	}{
		{
			`

Licensed(UDT+CP)/Effective      Udt     Sql     iPhtm   Pooled          Total

    ( 148 + 0   ) / 148         61      0       37      0               98

UDTNO USRNBR     UID  USRNAME USRTYPE              TTY         TIME        DATE
    1  17378     501   rkelly     udt            pts/2     08:21:57 Jan 25 2017
    2   5837     500  datatel     udt             udcs     00:01:02 Jan 15 2017
    3   5782     500  datatel     udt             udcs     00:01:01 Jan 15 2017
    4   2603     500  datatel     udt             udcs     09:20:48 Jan 25 2017

`,
			licenseStats{148, 0, 61, 0, 37, 0},
			true,
		},
	}

	for _, tc := range testCases {
		r := strings.NewReader(tc.input)
		s, err := parseLicenseStats(r)
		if (err != nil) == tc.ok {
			t.Fatalf("error mismatch: expected %v, got %v", tc.ok, err)
		}

		if !reflect.DeepEqual(*s, tc.want) {
			t.Error("state mismatch: expected %v, got %v", tc.want, *s)
		}
	}
}
