package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type licenseStats struct {
	limitUDT  float64
	limitCP   float64
	usageUDT  float64
	usageSQL  float64
	usageIPH  float64
	usagePool float64
}

func parseLicenseStats(r io.Reader) (*licenseStats, error) {
	var s *licenseStats

	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}

		v := strings.Fields(sc.Text())
		if len(v) != 12 && v[0] != "(" {
			continue
		}

		var err error
		s, err = newLicenseStats(v)
		if err != nil {
			return nil, err
		}

		break
	}
	if s == nil {
		return nil, errors.New("failed to parse license stats")
	}

	return s, nil
}

func newLicenseStats(v []string) (*licenseStats, error) {
	s := &licenseStats{}

	f, err := strconv.ParseFloat(v[1], 64)
	if err != nil {
		return nil, err
	}
	s.limitUDT = f

	f, err = strconv.ParseFloat(v[3], 64)
	if err != nil {
		return nil, err
	}
	s.limitCP = f

	f, err = strconv.ParseFloat(v[7], 64)
	if err != nil {
		return nil, err
	}
	s.usageUDT = f

	f, err = strconv.ParseFloat(v[8], 64)
	if err != nil {
		return nil, err
	}
	s.usageSQL = f

	f, err = strconv.ParseFloat(v[9], 64)
	if err != nil {
		return nil, err
	}
	s.usageIPH = f

	f, err = strconv.ParseFloat(v[10], 64)
	if err != nil {
		return nil, err
	}
	s.usagePool = f

	return s, nil
}
