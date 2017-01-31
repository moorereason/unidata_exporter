package main

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type unidataCollector struct {
	udtBinPath   string
	licenseUsage *prometheus.Desc
	licenseLimit *prometheus.Desc
	up           *prometheus.Desc
}

// newUnidataCollector returns a new unidataCollector using the given path to the
// UDTBIN.
func newUnidataCollector(udtbin string) *unidataCollector {
	return &unidataCollector{
		udtBinPath: udtbin,
		licenseUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "license", "usage"),
			"Unidata license usage.",
			[]string{"model", "src"},
			nil,
		),
		licenseLimit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "license", "limit"),
			"Unidata license limits.",
			[]string{"model"},
			nil,
		),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Unidata exporter success.",
			nil, nil,
		),
	}
}

// Describe sends the superset of all possible descriptors of metrics collected
// by this Collector to the provided channel and returns once the last
// descriptor has been sent.  See prometheus.Collector interface for more
// details.
func (c *unidataCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.licenseLimit
	ch <- c.licenseUsage
	ch <- c.up
}

// Collect sends each collected metric via the provided channel and returns once
// the last metric has been sent.  See prometheus.Collector interface for more
// details.
func (c *unidataCollector) Collect(ch chan<- prometheus.Metric) {
	var up float64

	err := c.CollectFromFile(filepath.Join(c.udtBinPath, "listuser"), ch)
	if err == nil {
		up = 1.0
	} else {
		log.Errorf("%v", err)
	}

	ch <- prometheus.MustNewConstMetric(
		c.up,
		prometheus.GaugeValue,
		up,
	)
}

func (c *unidataCollector) CollectFromReader(r io.Reader, ch chan<- prometheus.Metric) error {
	s, err := parseLicenseStats(r)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.licenseLimit,
		prometheus.GaugeValue,
		s.limitUDT,
		"udt",
	)

	ch <- prometheus.MustNewConstMetric(
		c.licenseLimit,
		prometheus.GaugeValue,
		s.limitCP,
		"cp",
	)

	ch <- prometheus.MustNewConstMetric(
		c.licenseUsage,
		prometheus.GaugeValue,
		s.usageUDT,
		"udt", "udt",
	)

	ch <- prometheus.MustNewConstMetric(
		c.licenseUsage,
		prometheus.GaugeValue,
		s.usageSQL,
		"udt", "sql",
	)

	ch <- prometheus.MustNewConstMetric(
		c.licenseUsage,
		prometheus.GaugeValue,
		s.usageIPH,
		"udt", "iphantom",
	)

	ch <- prometheus.MustNewConstMetric(
		c.licenseUsage,
		prometheus.GaugeValue,
		s.usagePool,
		"cp", "pool",
	)

	return nil
}

func (c *unidataCollector) CollectFromFile(path string, ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), *timeoutDuration)
	defer cancel()

	cmd := exec.CommandContext(ctx, path)

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	return c.CollectFromReader(bytes.NewReader(out), ch)
}

var _ prometheus.Collector = &unidataCollector{}
