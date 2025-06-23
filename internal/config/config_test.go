package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/goforj/godump"
	"github.com/google/go-cmp/cmp"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig("./testdata/invalid/hyprcircade.conf")
	if err == nil {
		t.Fatal("Error should be thrown")
	}

	expectedError := "ignoreAnchor property is missing in 1 files"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Expected error message: %s, got: %s", expectedError, err.Error())
	}

	cfg, err := NewConfig("./testdata/valid/hyprcircade.conf")
	if err != nil {
		t.Fatalf("Unexpected error was thrown: %s", err.Error())
	}

	expectedCfg := &Config{
		General: &General{
			Anchor:  "THEME_SWITCHER_TARGET",
			DarkAt:  20,
			LightAt: 8,
		},
		Files: []*File{
			&File{
				Path:         "./somefile.yaml",
				DayValue:     "light",
				NightValue:   "dark",
				IgnoreAnchor: false,
			},
			&File{
				Path:         "./somefile2.yaml",
				DayValue:     "light2",
				NightValue:   "dark2",
				IgnoreAnchor: true,
			},
		},
		Commands: []*Command{
			&Command{
				DayExec: "swww img some_light_image.png",
			},
			&Command{
				DayExec: "notify-send -t 3000 \"Switching theme to light\"",
			},
			&Command{
				NightExec: "swww img some_dark_image.png",
			},
			&Command{
				NightExec: "notify-send -t 3000 \"Switching theme to dark\"",
			},
		},
	}

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Fatalf(
			"Parsed configuration does not match expected one\nExpected: %v\nActual: %v\nDiff: %s",
			godump.DumpStr(expectedCfg),
			godump.DumpStr(cfg),
			cmp.Diff(expectedCfg, cfg),
		)
	}

	cfg, err = NewConfig("./testdata/valid-cmd/hyprcircade.conf")
	if err != nil {
		t.Fatalf("Unexpected error was thrown: %s", err.Error())
	}

	expectedCfg = &Config{
		General: &General{
			Anchor:  "THEME_SWITCHER_TARGET",
			DarkAt:  20,
			LightAt: 8,
		},
		Files: []*File{
			&File{
				Path:         "./somefile.yaml",
				DayValue:     "light",
				NightValue:   "dark",
				IgnoreAnchor: false,
			},
			&File{
				Path:         "./somefile2.yaml",
				DayValue:     "light2",
				NightValue:   "dark2",
				IgnoreAnchor: true,
			},
		},
		Commands: []*Command{
			&Command{
				DayExec: "swww img some_light_image.png",
			},
			&Command{
				DayExec: "notify-send -t 3000 \"Switching theme to light\"",
			},
			&Command{
				NightExec: "notify-send -t 3000 \"Switching theme to dark\"",
			},
			&Command{
				NightExec: "some night command",
			},
		},
	}

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Fatalf(
			"Parsed configuration does not match expected one\nExpected: %v\nActual: %v\nDiff: %s",
			godump.DumpStr(expectedCfg),
			godump.DumpStr(cfg),
			cmp.Diff(expectedCfg, cfg),
		)
	}
}
