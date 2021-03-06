// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"github.com/namsral/flag"

	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/datasync/resync"
	"github.com/ligato/cn-infra/db/keyval/etcdv3"
	"github.com/ligato/cn-infra/flavors/local"
)

// defines etcd & kafka flags // TODO switch to viper to avoid global configuration
func init() {
	flag.String("etcdv3-config", "etcd.conf",
		"Location of the Etcd configuration file; also set via 'ETCDV3_CONFIG' env variable.")
}

// FlavorEtcd glues together FlavorLocal plugins with ETCD & datasync plugin
// (which is useful for watching config.)
type FlavorEtcd struct {
	*local.FlavorLocal

	ETCD         etcdv3.Plugin
	ETCDDataSync kvdbsync.Plugin

	injected bool
}

// Inject sets object references
func (f *FlavorEtcd) Inject(resyncOrch *resync.Plugin) bool {
	if f.injected {
		return false
	}
	f.injected = true

	if f.FlavorLocal == nil {
		f.FlavorLocal = &local.FlavorLocal{}
	}
	f.FlavorLocal.Inject()

	f.ETCD.Deps.PluginInfraDeps = *f.InfraDeps("etcdv3")
	f.ETCDDataSync.Deps.PluginLogDeps = *f.LogDeps("etcdv3-datasync")
	f.ETCDDataSync.KvPlugin = &f.ETCD
	f.ETCDDataSync.ResyncOrch = resyncOrch
	f.ETCDDataSync.ServiceLabel = &f.ServiceLabel

	if f.StatusCheck.Transport == nil {
		f.StatusCheck.Transport = &f.ETCDDataSync
	}

	return true
}

// Plugins combines all Plugins in flavor to the list
func (f *FlavorEtcd) Plugins() []*core.NamedPlugin {
	f.Inject(nil)
	return core.ListPluginsInFlavor(f)
}
