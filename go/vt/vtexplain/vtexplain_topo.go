/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vtexplain

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"

	"github.com/youtube/vitess/go/vt/key"
	"github.com/youtube/vitess/go/vt/topo"
	"github.com/youtube/vitess/go/vt/vttablet/sandboxconn"

	topodatapb "github.com/youtube/vitess/go/vt/proto/topodata"
	vschemapb "github.com/youtube/vitess/go/vt/proto/vschema"
)

// VTExplainTopo satisfies the SrvTopoServer interface.
// modeled after the vtgate test sandboxTopo
type VTExplainTopo struct {
	// Map of keyspace name to vschema
	Keyspaces map[string]*vschemapb.Keyspace

	// Map of ks/shard to test tablet connection
	TabletConns map[string]*sandboxconn.SandboxConn

	// Synchronization lock
	Lock sync.Mutex
}

func (et *VTExplainTopo) getSrvVSchema() *vschemapb.SrvVSchema {
	et.Lock.Lock()
	defer et.Lock.Unlock()

	return &vschemapb.SrvVSchema{
		Keyspaces: et.Keyspaces,
	}
}

// GetSrvKeyspaceNames is part of SrvTopoServer.
func (et *VTExplainTopo) GetSrvKeyspaceNames(ctx context.Context, cell string) ([]string, error) {
	et.Lock.Lock()
	defer et.Lock.Unlock()

	keyspaces := make([]string, 0, 1)
	for k := range et.Keyspaces {
		keyspaces = append(keyspaces, k)
	}
	return keyspaces, nil
}

// GetSrvKeyspace is part of SrvTopoServer.
func (et *VTExplainTopo) GetSrvKeyspace(ctx context.Context, cell, keyspace string) (*topodatapb.SrvKeyspace, error) {
	et.Lock.Lock()
	defer et.Lock.Unlock()

	vschema := et.Keyspaces[keyspace]
	if vschema == nil {
		return nil, fmt.Errorf("no vschema for keyspace %s", keyspace)
	}

	if vschema.Sharded {
		shards := make([]*topodatapb.ShardReference, 0, NUM_SHARDS)
		for i := 0; i < NUM_SHARDS; i++ {
			kr, err := key.EvenShardsKeyRange(i, NUM_SHARDS)
			if err != nil {
				return nil, err
			}

			shard := &topodatapb.ShardReference{
				Name:     key.KeyRangeString(kr),
				KeyRange: kr,
			}
			shards = append(shards, shard)
		}

		shardedSrvKeyspace := &topodatapb.SrvKeyspace{
			ShardingColumnName: "", // exact value is ignored
			ShardingColumnType: 0,
			Partitions: []*topodatapb.SrvKeyspace_KeyspacePartition{
				{
					ServedType:      topodatapb.TabletType_MASTER,
					ShardReferences: shards,
				},
				{
					ServedType:      topodatapb.TabletType_REPLICA,
					ShardReferences: shards,
				},
				{
					ServedType:      topodatapb.TabletType_RDONLY,
					ShardReferences: shards,
				},
			},
		}
		return shardedSrvKeyspace, nil

	} else {
		// unsharded
		kr, err := key.EvenShardsKeyRange(0, 1)
		if err != nil {
			return nil, err
		}

		shard := &topodatapb.ShardReference{
			Name: key.KeyRangeString(kr),
		}

		unshardedSrvKeyspace := &topodatapb.SrvKeyspace{
			Partitions: []*topodatapb.SrvKeyspace_KeyspacePartition{
				{
					ServedType:      topodatapb.TabletType_MASTER,
					ShardReferences: []*topodatapb.ShardReference{shard},
				},
				{
					ServedType:      topodatapb.TabletType_REPLICA,
					ShardReferences: []*topodatapb.ShardReference{shard},
				},
				{
					ServedType:      topodatapb.TabletType_RDONLY,
					ShardReferences: []*topodatapb.ShardReference{shard},
				},
			},
		}

		return unshardedSrvKeyspace, nil
	}
}

// WatchSrvVSchema is part of SrvTopoServer.
func (et *VTExplainTopo) WatchSrvVSchema(ctx context.Context, cell string) (*topo.WatchSrvVSchemaData, <-chan *topo.WatchSrvVSchemaData, topo.CancelFunc) {
	return &topo.WatchSrvVSchemaData{
		Value: et.getSrvVSchema(),
	}, make(chan *topo.WatchSrvVSchemaData), func() {}
}
