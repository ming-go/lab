/*
	Ref:
		[1]: https://cloud.tencent.com/developer/article/1183490
		[2]: https://github.com/rqlite/rqlite
		[3]: https://github.com/yongman/leto
		[4]: https://github.com/hashicorp/consul/tree/master/agent/consul
*/

package hraft

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/consul/agent/consul/fsm"
	"github.com/hashicorp/consul/agent/consul/state"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"go.uber.org/zap"
)

type Config struct {
	DevMode   bool
	Bootstrap bool
	LogOutput io.Writer
	LogLevel  string
}

type hraft struct {
	Raft          *raft.Raft
	fsm           *fsm.FSM
	tombstoneGC   *state.TombstoneGC
	logOutput     io.Writer //raft.LogStore
	log           raft.LogStore
	logger        *log.Logger
	stable        raft.StableStore
	sbao          raft.SnapshotStore
	RaftConfig    *raft.Config
	raftStore     *raftboltdb.BoltStore
	raftNotifyCh  chan bool
	raftInmem     *raft.InmemStore
	raftTransport *raft.NetworkTransport
	DevMode       bool
	Bootstrap     bool
}

var TombstoneTTL time.Duration = 15 * time.Minute
var TombstoneTTLGranularity time.Duration = 30 * time.Second
var raftLogCacheSize = 512
var snapshotsRetained = 2

const (
	raftState = "raft/"
)

func ensurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

func New(nodeID raft.ServerID, bind string, bootstrap bool) (*hraft, error) {
	var err error
	hraft := hraft{
		RaftConfig: raft.DefaultConfig(),
		Bootstrap:  bootstrap,
	}

	hraft.RaftConfig.LogLevel = "INFO"

	hraft.tombstoneGC, err = state.NewTombstoneGC(
		TombstoneTTL,
		TombstoneTTLGranularity,
	)
	if err != nil {
		return nil, err
	}

	hraft.logOutput = zap.NewStdLog(zap.L()).Writer()

	// TODO: Replace by LevelDB or Rocksdb or BoltDB or Badger
	// Create the FSM
	hraft.fsm, err = fsm.New(hraft.tombstoneGC, hraft.logOutput)
	if err != nil {
		return nil, err
	}

	hraft.RaftConfig.LocalID = nodeID
	hraft.logger = log.New(os.Stderr, "", log.LstdFlags)

	addr, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		return nil, err
	}

	// Create a transport layer.
	// TODO:
	//transConfig := &raft.NetworkTransportConfig{
	//	MaxPool: 3,
	//	Timeout: 10 * time.Second,
	//	Logger:  hraft.logger,
	//}

	//trans := raft.NewNetworkTransportWithConfig(transConfig)
	trans, err := raft.NewTCPTransport(bind, addr, 3, 10*time.Second, hraft.logOutput)
	if err != nil {
		return nil, err
	}
	hraft.raftTransport = trans

	var log raft.LogStore
	var stable raft.StableStore
	var snap raft.SnapshotStore

	if hraft.DevMode {
		store := raft.NewInmemStore()
		hraft.raftInmem = store
		stable = store
		log = store
		snap = raft.NewInmemSnapshotStore()
	} else {
		raftDataPath := filepath.Join("./data", raftState)
		if err := ensurePath(raftDataPath, true); err != nil {
			return nil, err
		}

		// Create the backend raft store for logs and stable storage.
		store, err := raftboltdb.NewBoltStore(filepath.Join(raftDataPath, "raft.db"))
		if err != nil {
			return nil, err
		}

		hraft.raftStore = store
		stable = store

		// Wrap the store in a LogCache to improve performance.
		cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
		if err != nil {
			return nil, err
		}
		log = cacheStore

		snapshots, err := raft.NewFileSnapshotStore(raftDataPath, snapshotsRetained, hraft.logOutput)
		if err != nil {
			return nil, err
		}
		snap = snapshots

		peersFile := filepath.Join(raftDataPath, "peers.json")
		peersInfoFile := filepath.Join(raftDataPath, "peers.info")

		if _, err := os.Stat(peersInfoFile); os.IsNotExist(err) {
			if err := ioutil.WriteFile(peersInfoFile, []byte(peersInfoContent), 0755); err != nil {
				return nil, fmt.Errorf("failed to write peers.info file: %v", err)
			}

			if _, err := os.Stat(peersFile); err == nil {
				if err := os.Remove(peersFile); err != nil {
					return nil, fmt.Errorf("failed to delete peers.json, please delete manually (see peers.info for details): %v", err)
				}

				zap.L().Info("consul: deleted peers.json file (see peers.info for details)")
			}
		} else if _, err := os.Stat(peersFile); err == nil {
			zap.L().Info("consul: found peers.json file, recovering Raft configuration...")

			var configuration raft.Configuration
			configuration, err = raft.ReadConfigJSON(peersFile)
			if err != nil {
				return nil, fmt.Errorf("recovery failed to parse peers.json: %v", err)
			}

			tmpFsm, err := fsm.New(hraft.tombstoneGC, hraft.logOutput)
			if err != nil {
				return nil, fmt.Errorf("recovery failed to make temp FSM: %v", err)
			}

			if err := raft.RecoverCluster(hraft.RaftConfig, tmpFsm, log, stable, snap, trans, configuration); err != nil {
				return nil, fmt.Errorf("recovery failed: %v", err)
			}

			if err := os.Remove(peersFile); err != nil {
				return nil, fmt.Errorf("recovery failed to delete peers.json, please delete manually (see peers.info for details): %v", err)
			}

			zap.L().Info("consul: deleted peers.json file after successful recovery")
		}
	}

	if hraft.Bootstrap || hraft.DevMode {
		hasState, err := raft.HasExistingState(log, stable, snap)
		if err != nil {
			return nil, err
		}

		if !hasState {
			configuration := raft.Configuration{
				Servers: []raft.Server{
					raft.Server{
						ID:      hraft.RaftConfig.LocalID,
						Address: trans.LocalAddr(),
					},
				},
			}
			if err := raft.BootstrapCluster(hraft.RaftConfig, log, stable, snap, trans, configuration); err != nil {
				return nil, err
			}
		}
	}

	raftNotifyCh := make(chan bool, 1)
	hraft.RaftConfig.NotifyCh = raftNotifyCh
	hraft.raftNotifyCh = raftNotifyCh

	// Setup the Raft store.
	hraft.Raft, err = raft.NewRaft(hraft.RaftConfig, hraft.fsm, log, stable, snap, trans)
	if err != nil {
		return nil, err
	}

	return &hraft, err
}

func (hr *hraft) Join(nodeID string, addr string) error {
	hr.logger.Printf("received join request for remote node %s, addr %s", nodeID, addr)

	cf := hr.Raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		hr.logger.Printf("failed to get raft configuration")
		return err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(nodeID) {
			hr.logger.Printf("node %s already joined raft cluster", nodeID)
			return nil
		}
	}

	f := hr.Raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if err := f.Error(); err != nil {
		return err
	}

	hr.logger.Printf("node %s at %s joined successfully", nodeID, addr)

	return nil
}
