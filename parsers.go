package ezcli

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

func parseDurationSlice(s string) ([]time.Duration, error) {
	if len(s) == 0 {
		return []time.Duration{}, nil
	}
	// Flag format of: "[durationString,durationString]"
	// Env format of: "durationString durationString"

	// Remove any containing brackets if they exist
	if s[0] == []byte("[")[0] && s[len(s)-1] == []byte("]")[0] {
		s = s[1 : len(s)-1]
	}

	// Handle splitting on different types of values
	// Commas used to seperate each duration value
	durationStrings := strings.Split(s, ",")
	// If we didn't get more than one value by splitting on commas, try spaces
	if len(durationStrings) == 1 {
		durationStrings = strings.Split(s, " ")
	}

	durations := make([]time.Duration, len(durationStrings))
	for i, durationString := range durationStrings {
		d, err := time.ParseDuration(durationString)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse from %s", s)
		}
		durations[i] = d
	}
	return durations, nil
}
