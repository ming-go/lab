package hraft

import (
	"io"
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

	path := filepath.Join("./data", raftState)
	if err := ensurePath(path, true); err != nil {
		return err
	}

	store, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft.db"))
	if err != nil {
		return nil, err
	}

	hraft.raftStore = store

	return nil, nil
}
