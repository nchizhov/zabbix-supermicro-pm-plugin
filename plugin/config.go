package zbplugin

import (
	"os"
	"support"
	"time"

	"golang.zabbix.com/sdk/conf"
	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/plugin"
)

const (
	discovery = "smpm.discovery"
	fieldGet  = "smpm.get"
)

type session struct {
	User     string `conf:"optional"`
	Password string `conf:"optional"`
}

type PluginConfig struct {
	System          plugin.SystemOptions `conf:"name=System,optional"`
	IPMITool        string               `conf:"name=IPMITool"`
	CollectInterval time.Duration        `conf:"name=CollectInterval,optional,range=1:30,default=1"`

	Sessions map[string]session `conf:"optional"`
}

// Take config
func (p *smpm) Configure(global *plugin.GlobalOptions, option any) {
	pConfig := &PluginConfig{}

	if err := conf.UnmarshalStrict(option, pConfig); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err.Error())
		return
	}

	p.config = pConfig
}

// Validate config
func (p *smpm) Validate(options any) error {
	var opts PluginConfig

	if err := conf.UnmarshalStrict(options, &opts); err != nil {
		return errs.Wrap(err, "failed to unmarshal configuration options")
	}

	if opts.IPMITool == "" {
		return errs.Errorf(
			"IPMITool cannot be empty",
		)
	}
	IPMIToolInfo, err := os.Stat(opts.IPMITool)
	if err != nil {
		return errs.Errorf(
			"IPMITool by path %s not exists",
			opts.IPMITool,
		)
	}
	if !support.IsExecutable(opts.IPMITool, IPMIToolInfo) {
		return errs.Errorf(
			"IPMITool by path %s not executable",
			opts.IPMITool,
		)
	}

	return nil
}

func (p *smpm) registerMetrics() error {
	p.metrics = map[string]*smpmMetric{
		discovery: {
			Description: "Discovery SuperMicro IPMI Power Supply Modules.",
			Params:      0,
			Handler:     "GetDiscovery",
		},
		fieldGet: {
			Description: "Get SuperMicro IPMI Power Supply Module Param.",
			Params:      2,
			Handler:     "GetFieldData",
		},
	}

	metrics := []string{}

	for k, m := range p.metrics {
		metrics = append(metrics, k, m.Description)
	}

	err := plugin.RegisterMetrics(
		p,
		Name,
		metrics...,
	)
	if err != nil {
		return errs.Wrap(err, "failed to register metrics")
	}

	return nil
}
