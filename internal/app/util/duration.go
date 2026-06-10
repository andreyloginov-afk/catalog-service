package util

import "time"

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	d.Duration = duration

	return nil
}
