package tools

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type node struct {
	lock      sync.Mutex
	nodeBits  uint8
	seqBits   uint8
	timeShift uint8
	seqMask   int64
	epoch     int64
	timeStamp int64
	nodeId    int64
	seq       int64
}

type ID int64

const (
	defaultNodeBits uint8 = 10
	defaultSeqBits  uint8 = 12
	defaultNodeId   int64 = 1<<5 | 1
	defaultEpoch    int64 = 1561392000000
)

func newDefaultNode() *node {
	n, _ := newNode(defaultNodeBits, defaultSeqBits, defaultEpoch, defaultNodeId)
	return n
}

func newNode(nodeBits, seqBits uint8, epoch, nodeId int64) (*node, error) {
	if nodeBits+seqBits >= 63 {
		return nil, errors.New("nodeBits and seqBits too long")
	}
	return &node{
		seqBits:   seqBits,
		seqMask:   -1 ^ (-1 << seqBits),
		nodeBits:  nodeBits,
		timeShift: seqBits + nodeBits,
		epoch:     epoch,
		nodeId:    nodeId,
	}, nil
}

func (n *node) ID() ID {
	n.lock.Lock()
	defer n.lock.Unlock()
	now := time.Now().UnixNano() / 1e6
	if now == n.timeStamp {
		n.seq = (n.seq + 1) & n.seqMask
		if n.seq == 0 {
			for now <= n.timeStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		n.seq = 0
	}
	n.timeStamp = now
	return ID((now-n.epoch)<<n.timeShift | n.nodeId<<n.seqBits | n.seq)
}

var _node = newDefaultNode()

func GenerateID() ID {
	return _node.ID()
}

func (i ID) Int64() int64 {
	return int64(i)
}

func (i ID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i ID) Base36() string {
	return strconv.FormatInt(int64(i), 36)
}
