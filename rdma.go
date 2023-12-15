package rdma

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus-community/windows_exporter/pkg/perflib"
	"github.com/prometheus-community/windows_exporter/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Name           = "rdma"
	FlagNicExclude = "collector.net.nic-exclude"
	FlagNicInclude = "collector.net.nic-include"

	FlagRdmaExclude = "collector.net.nic-exclude"
	FlagRdmaInclude = "collector.net.nic-include"
)

type Config struct {
	NicInclude string `yaml:"nic_include"`
	NicExclude string `yaml:"nic_exclude"`
}

//	type Config struct {
//		RdmaInclude string `yaml:"rdma_include"`
//		RdmaExclude string `yaml:"rdma_exclude"`
//	}
var ConfigDefaults = Config{
	NicInclude: ".+",
	NicExclude: "",
}

//	var ConfigDefaults = Config{
//		RdmaInclude: ".+",
//		RdmaExclude: "",
//	}
var nicNameToUnderscore = regexp.MustCompile("[^a-zA-Z0-9]")

//var rdmaNameToUnderscore = regexp.MustCompile("[^a-zA-Z0-9]")

// A collector is a Prometheus collector for RDMA metrics
type collector struct {
	logger log.Logger

	//rdmaInclude *string
	//rdmaExclude *string
	nicInclude *string
	nicExclude *string

	// Modify metric names based on RDMA counters
	// ...

	RDMAInboundBytes  *prometheus.Desc
	RDMAOutboundBytes *prometheus.Desc

	//rdmaIncludePattern *regexp.Regexp
	//rdmaExcludePattern *regexp.Regexp
	nicIncludePattern *regexp.Regexp
	nicExcludePattern *regexp.Regexp
}

// New initializes a new instance of the RDMA collector.
func New(logger log.Logger, config *Config) types.Collector {
	if config == nil {
		config = &ConfigDefaults
	}

	c := &collector{
		nicExclude: &config.NicExclude,
		nicInclude: &config.NicInclude,
		//rdmaExclude: &config.RdmaExclude,
		//rdmaInclude: &config.RdmaInclude,
	}
	c.SetLogger(logger)
	return c
}

// NewWithFlags initializes a new RDMA collector with flags for configuration.
func NewWithFlags(app *kingpin.Application) types.Collector {
	c := &collector{
		nicInclude: app.Flag(
			FlagNicInclude,
			"Regexp of NIC:s to include. NIC name must both match include and not match exclude to be included.",
		).Default(ConfigDefaults.NicInclude).String(),

		nicExclude: app.Flag(
			FlagNicExclude,
			"Regexp of NIC:s to exclude. NIC name must both match include and not match exclude to be included.",
		).Default(ConfigDefaults.NicExclude).String(),
	}

	return c

	//c := &collector{
	//	rdmaInclude: app.Flag(
	//		FlagRdmaInclude,
	//		"Regexp of RDMA interfaces to include. Interface name must both match include and not match exclude to be included.",
	//	).Default(ConfigDefaults.RdmaInclude).String(),

	//	rdmaExclude: app.Flag(
	//		FlagRdmaExclude,
	//		"Regexp of RDMA interfaces to exclude. Interface name must both match include and not match exclude to be included.",
	//	).Default(ConfigDefaults.RdmaExclude).String(),
	//}

	//return c
}

func (c *collector) GetName() string {
	return Name
}

func (c *collector) SetLogger(logger log.Logger) {
	c.logger = log.With(logger, "collector", Name)
}
func (c *collector) GetPerfCounter() ([]string, error) {
	return []string{"RDMA Activity"}, nil
}

// Build initializes RDMA metrics and compiles regular expression patterns for inclusion and exclusion.
func (c *collector) Build() error {
	// Modify metric names based on RDMA counters
	// ...
	c.RDMAInboundBytes = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "RDMA_inbound_bytes"),
		"(RDMA.RDMAInboundBytesPerSec)",
		[]string{"nic"},
		nil,
	)

	// Add additional metric initialization based on your RDMA technology
	// ...
	var err error
	c.nicIncludePattern, err = regexp.Compile(fmt.Sprintf("^(?:%s)$", *c.nicInclude))
	if err != nil {
		return err
	}

	c.nicExcludePattern, err = regexp.Compile(fmt.Sprintf("^(?:%s)$", *c.nicExclude))
	if err != nil {
		return err
	}

	return nil
	//var err error
	//c.rdmaIncludePattern, err = regexp.Compile(fmt.Sprintf("^(?:%s)$", *c.rdmaInclude))
	//if err != nil {
	//	return err
	//}
	//c.rdmaExcludePattern, err = regexp.Compile(fmt.Sprintf("^(?:%s)$", *c.rdmaExclude))
	//if err != nil {
	//	return err
	//}
	//return nil
}

// Collect collects RDMA metrics and sends them to the Prometheus metric channel.
func (c *collector) Collect(ctx *types.ScrapeContext, ch chan<- prometheus.Metric) error {
	// Collect RDMA metrics
	// ...
	if desc, err := c.collect(ctx, ch); err != nil {
		_ = level.Error(c.logger).Log("failed collecting RDMA metrics", "desc", desc, "err", err)
		return err
	}
	return nil
	// Add additional metric collections based on your RDMA technology
	// ...
}

func mangleNetworkName(name string) string {
	return nicNameToUnderscore.ReplaceAllString(name, "_")
}

// Helper function to modify RDMA interface name
//func mangleRdmaName(name string) string {
//	return rdmaNameToUnderscore.ReplaceAllString(name, "_")
//}

// RDMA interface struct
type rdmaInterface struct {
	// Modify fields based on RDMA counters
	// ...
	RDMAInboundBytesPerSec float64 `perflib:"RDMA Inbound Bytes/sec"`
	Name                   string

	// Add additional counters specific to your RDMA technology
	// ...

}

// Collection function to collect RDMA metrics
func (c *collector) collect(ctx *types.ScrapeContext, ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []rdmaInterface

	if err := perflib.UnmarshalObject(ctx.PerfObjects["RDMA Activity"], &dst, c.logger); err != nil {
		return nil, err
	}
	for _, nic := range dst {
		if c.nicExcludePattern.MatchString(nic.Name) ||
			!c.nicIncludePattern.MatchString(nic.Name) {
			continue
		}
		//for _, rdma := range dst {
		//	if c.rdmaExcludePattern.MatchString(rdma.Name) ||
		//		!c.rdmaIncludePattern.MatchString(rdma.Name) {
		//		continue
		//	}
		name := mangleNetworkName(nic.Name)
		if name == "" {
			continue
		}

		//	name := mangleRdmaName(rdma.Name)
		//	if name == "" {
		//		continue
		//	}

		// Collect RDMA metrics
		// ...
		ch <- prometheus.MustNewConstMetric(
			c.RDMAInboundBytes,
			prometheus.CounterValue,
			nic.RDMAInboundBytesPerSec,
			name,
		)

	}
	return nil, nil
}
