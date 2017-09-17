package ui

import "testing"

func TestGetHeight(t *testing.T) {
	cases := []struct {
		input          string
		expectedHeight int
	}{
		{
			``, 0,
		},
		{
			`Single Line`, 1,
		},
		{
			`__    __     __  __     __         ______   __     ______   __         ______
			/\ "-./  \   /\ \/\ \   /\ \       /\__  _\ /\ \   /\  == \ /\ \       /\  ___\
			\ \ \-./\ \  \ \ \_\ \  \ \ \____  \/_/\ \/ \ \ \  \ \  _-/ \ \ \____  \ \  __\
			 \ \_\ \ \_\  \ \_____\  \ \_____\    \ \_\  \ \_\  \ \_\    \ \_____\  \ \_____\
			  \/_/  \/_/   \/_____/   \/_____/     \/_/   \/_/   \/_/     \/_____/   \/_____/`, 5,
		},
	}

	for _, c := range cases {
		height := getHeight(c.input)
		if height != c.expectedHeight {
			t.Errorf("Calculated incorrect height, got %d, expected %d", height, c.expectedHeight)
		}
	}
}

func BenchmarkGetHeight(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getHeight(
			`__    __     __  __     __         ______   __     ______   __         ______
			/\ "-./  \   /\ \/\ \   /\ \       /\__  _\ /\ \   /\  == \ /\ \       /\  ___\
			\ \ \-./\ \  \ \ \_\ \  \ \ \____  \/_/\ \/ \ \ \  \ \  _-/ \ \ \____  \ \  __\
			 \ \_\ \ \_\  \ \_____\  \ \_____\    \ \_\  \ \_\  \ \_\    \ \_____\  \ \_____\
			  \/_/  \/_/   \/_____/   \/_____/     \/_/   \/_/   \/_/     \/_____/   \/_____/`,
		)
	}
}

func TestGetWidth(t *testing.T) {
	cases := []struct {
		input         string
		expectedWidth int
	}{
		{
			``, 0,
		},
		{
			`-----`, 5,
		},
		{
			`
------
---------
----
----------------------
---------
			`, 22,
		},
	}

	for _, c := range cases {
		width := getWidth(c.input)
		if width != c.expectedWidth {
			t.Errorf("Calculated incorrect width, got %d, expected %d", width, c.expectedWidth)
		}
	}
}

func BenchmarkGetWidth(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getWidth(
			`
----
----------
--
-------
-------------------
------
			`,
		)
	}
}
