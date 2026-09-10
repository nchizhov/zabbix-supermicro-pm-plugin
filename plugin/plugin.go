package zbplugin

import (
	"context"
	"ipmi"
	"support"
	"time"

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/log"
	"golang.zabbix.com/sdk/plugin/container"

	"golang.zabbix.com/sdk/plugin"
)

const (
	Name   = "SMIPMIps"
	Url    = "https://api.github.com/repos/nchizhov/zabbix-supermicro-pm-plugin/releases/latest"
	Prefix = "smpm"
)

type smpm struct {
	plugin.Base

	config         *PluginConfig
	metrics        map[string]*smpmMetric
	IPMICollect    chan struct{}
	IPMICollectCtx context.Context
	IPMIData       ipmi.Data
}

type smpmMetric struct {
	Description string
	Params      int
	Handler     string
}

func New(ctx context.Context) (*smpm, error) {
	p := &smpm{}

	err := log.Open(log.Console, log.Info, "", 0)
	if err != nil {
		return nil, errs.Wrap(err, "failed to open log")
	}

	p.Logger = log.New(Name)
	p.IPMICollectCtx = ctx

	err = p.registerMetrics()
	if err != nil {
		return nil, errs.Wrap(err, "plugin failed to register metrics")
	}

	p.IPMICollect = make(chan struct{})

	return p, nil
}

// Start
func (p *smpm) Start() {
	isNewUpdate, updateInfo := support.CheckUpdate(Url, Prefix)
	if isNewUpdate {
		p.Infof(updateInfo)
	}
	p.collectInfo()
	go func() {
		ticker := time.NewTicker(p.config.CollectInterval * time.Minute)
		defer ticker.Stop()

		p.Infof("collector started")

		for {
			select {
			case <-ticker.C:
				p.collectInfo()
			case <-p.IPMICollect:
				p.Infof("collector stopped")
				return
			}
		}
	}()

	p.Infof("start plugin %s", Name)
}

// Stop
func (p *smpm) Stop() {
	close(p.IPMICollect)
	p.Infof("stop plugin %s", Name)
}

// Run
func (p *smpm) Run() error {
	h, err := container.NewHandler(Name)
	if err != nil {
		return errs.Wrap(err, "failed to create new handler")
	}

	p.Logger = h

	err = h.Execute()
	if err != nil {
		return errs.Wrap(err, "failed to execute plugin handler")
	}

	return nil
}

func (p *smpm) Export(
	key string, params []string, ctx plugin.ContextProvider,
) (any, error) {
	if metric, ok := p.metrics[key]; ok {
		totalParams := len(params)
		if totalParams != metric.Params {
			return nil, errs.Errorf("incorrect number of params. Required %d", metric.Params)
		}
		data, err := ipmi.ReflectHandler(&p.IPMIData, params, metric.Handler)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return data, nil
		}
		if totalParams == 2 {
			return nil, errs.Errorf("unsupported param %q", params[1])
		}
		return nil, errs.Errorf("undefined data for key %q", key)
	}
	return nil, errs.Errorf("unknown metric key %q", key)
}

func (p *smpm) collectInfo() {
	pmData, err := ipmi.GetInfo(p.IPMICollectCtx, p.config.IPMITool)
	if err != nil {
		p.Infof("collect info error: %s", err.Error())
	} else {
		p.IPMIData.Update(pmData)
	}
}
