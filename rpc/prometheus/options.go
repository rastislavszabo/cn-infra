package prometheus

import (
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logrus"
	"github.com/ligato/cn-infra/rpc/rest"
)

// DefaultPlugin is default instance of Plugin
var DefaultPlugin = NewPlugin()

// NewPlugin creates a new Plugin with the provides Options
func NewPlugin(opts ...Option) *Plugin {
	p := &Plugin{}

	for _, o := range opts {
		o(p)
	}

	if p.Deps.PluginName == "" {
		p.Deps.PluginName = "prometheus"
	}
	if p.Deps.Log == nil {
		p.Deps.Log = logging.ForPlugin(p.Deps.PluginName.String(), logrus.DefaultRegistry)
	}
	if p.Deps.HTTP == nil {
		p.Deps.HTTP = rest.DefaultPlugin
	}

	return p
}

// Option is a function that acts on a Plugin to inject Dependencies or configuration
type Option func(*Plugin)

// UseDeps injects a particular set of Dependencies
func UseDeps(deps Deps) Option {
	return func(p *Plugin) {
		p.Deps = deps
	}
}
