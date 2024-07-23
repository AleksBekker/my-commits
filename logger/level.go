package logger

import (
	"github.com/AleksBekker/my-commits/colors"
	"strings"
)

type Level struct {
	prefix       string
	minVerbosity int
	clrs         colors.Colors
}

func DefaultLevel(tag string) Level {
	settings, ok := map[string]struct {
		prefix string
		verb   int
		clrs   colors.Colors
	}{
		"error": {"ERROR: ", Err, []string{colors.RedFg}},
		"fatal": {"FATAL: ", None, []string{colors.RedBg, colors.BlackFg}},
		"info":  {"INFO: ", Info, []string{colors.CyanFg}},
		"panic": {"PANIC: ", None, []string{colors.RedBg, colors.BlackFg}},
		"warn":  {"WARNING: ", Warn, []string{colors.YellowFg}},
	}[tag]
	prefix, verb, clrs := settings.prefix, settings.verb, settings.clrs

	if !ok {
		prefix = strings.ToUpper(tag) + ": "
		verb, clrs = All, nil
	}

	return Level{prefix, verb, clrs}
}
