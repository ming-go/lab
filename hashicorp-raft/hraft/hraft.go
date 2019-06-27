package hraft

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/consul/agent/consul/fsm"
	"github.com/hashicorp/consul/agent/consul/state"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"go.uber.org/zap"
)

type hraft struct {
	raft        *raft.Raft
	fsm         *fsm.FSM
	tombstoneGC *state.TombstoneGC
	logOutput   io.Writer //raft.LogStore
	log         raft.LogStore
	stable      raft.StableStore
	sbao        raft.SnapshotStore
	RaftConfig  *raft.Config
	raftStore   *raftboltdb.BoltStore
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

func New(nodeID raft.ServerID) (*hraft, error) {
	var err error
	hraft := hraft{
		RaftConfig: raft.DefaultConfig(),
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
	hraft.fsm, err = fsm.New(hraft.tombstoneGC, hraft.logOutput)
	if err != nil {
		return nil, err
	}

	hraft.RaftConfig.LocalID = nodeID

	var log raft.LogStore
	var stable raft.StableStore
	var snap raft.SnapshotStore

	raftDataPath := filepath.Join("./data", raftState)
	if err := ensurePath(path, true); err != nil {
		return err
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

	snapshots, err := raft.NewFileSnapshotStore(filepath, hraft.logOutput)
	if err != nil {
		return err
	}
	snap = snapshots

	peersFile := filepath.Join(path, "peers.json")
	peersInfoFile := filepath.Join(path, "peers.info")

	if _, err := os.Stat(peersInfoFile); os.IsNotExist(err) {
		if err := ioutil.WriteFile(peersInfoFile, []byte(peersInfoContent), 0755); err != nil {
			return fmt.Errorf("failed to write peers.info file: %v", err)
		}

		if _, err := os.Stat(peersFile); err == nil {
			if err := os.Remove(peersFile); err != nil {
				return fmt.Errorf("failed to delete peers.json, please delete manually (see peers.info for details): %v", err)
			}

			zap.L().Info("consul: deleted peers.json file (see peers.info for details)")
		}
	} else if _, err := os.Stat(peersFile); err == nil {
		zap.L().Info("consul: found peers.json file, recovering Raft configuration...")

		var configuration raft.Configuration
		configuration, err = raft.ReadConfigJSON(peersFile)
		if err != nil {
			return fmt.Errorf("recovery failed to parse peers.json: %v", err)
		}

		tmpFsm, err := fsm.New(hraft.tombstoneGC, hraft.LogOutput)
		if err != nil {
			return fmt.Errorf("recovery failed to make temp FSM: %v", err)
		}

	}

	return nil, nil
}
