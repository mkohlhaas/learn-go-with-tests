package main

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         any
		ExpectedCalls []string
	}{
		{
			"struct with one string field",
			struct{ Name string }{"Chris"},
			[]string{"Chris"},
		},
		{
			"struct with two string fields",
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},
		{
			"struct with non string field",
			struct {
				Name string
				Age  int
			}{"Chris", 33},
			[]string{"Chris"},
		},
		{
			"nested fields",
			Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"pointers to things",
			&Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"slices",
			[]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},
		{
			"arrays",
			[2]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := []string{}

			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			assertEqual(t, got, test.ExpectedCalls)
		})
	}

	t.Run("with maps", func(t *testing.T) {
		got := []string{}

		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}

		walk(aMap, func(input string) {
			got = append(got, input)
		})

		if len(aMap) != len(got) {
			t.Errorf("map %v has length %v instead of wanted length %v", aMap, len(aMap), len(got))
		}
		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})

	t.Run("with channels", func(t *testing.T) {
		got := []string{}
		want := []string{"Berlin", "Katowice"}

		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{33, "Berlin"}
			aChannel <- Profile{34, "Katowice"}
			close(aChannel)
		}()

		walk(aChannel, func(input string) {
			got = append(got, input)
		})

		assertEqual(t, got, want)
	})

	t.Run("with function", func(t *testing.T) {
		got := []string{}
		want := []string{"Berlin", "Katowice"}

		aFunction := func() (Profile, Profile) {
			return Profile{33, "Berlin"}, Profile{34, "Katowice"}
		}

		walk(aFunction, func(input string) {
			got = append(got, input)
		})

		assertEqual(t, got, want)
	})

	t.Run("with side effects only function", func(t *testing.T) {
		got := []string{}
		want := []string{}

		aFunction := func() {}

		walk(aFunction, func(input string) {
			got = append(got, input)
		})

		assertEqual(t, got, want)
	})
}

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func assertContains(t testing.TB, haystack []string, needle string) {
	t.Helper()
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
			break
		}
	}
	if !contains {
		t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
	}
}

func assertEqual(t *testing.T, got, want []string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
