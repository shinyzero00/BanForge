package blocker

import (
	"fmt"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

// BlockerType defines the type of firewall blocker
type BlockerType string

const (
	BlockerTypeNftables  BlockerType = "nftables"
	BlockerTypeIptables  BlockerType = "iptables"
	BlockerTypeFirewalld BlockerType = "firewalld"
	BlockerTypeUfw       BlockerType = "ufw"
)

// BlockerFactory creates new blocker instances
type BlockerFactory struct {
	logger *logger.Logger
}

// NewBlockerFactory creates a new blocker factory
func NewBlockerFactory(logger *logger.Logger) *BlockerFactory {
	return &BlockerFactory{
		logger: logger,
	}
}

// Create creates a new blocker instance of the specified type
func (bf *BlockerFactory) Create(btype BlockerType, config string) (BlockerEngine, error) {
	switch btype {
	case BlockerTypeNftables:
		return NewNftables(bf.logger, config), nil
	case BlockerTypeIptables:
		return NewIptables(bf.logger, config), nil
	case BlockerTypeFirewalld:
		return NewFirewalld(bf.logger), nil
	case BlockerTypeUfw:
		return NewUfw(bf.logger), nil
	default:
		return nil, fmt.Errorf("unknown blocker type: %s", btype)
	}
}

// CreateFromString creates a blocker from string type name
func (bf *BlockerFactory) CreateFromString(typename, config string) (BlockerEngine, error) {
	return bf.Create(BlockerType(typename), config)
}

// ListAvailable returns all available blocker types
func ListAvailable(logger *logger.Logger) []string {
	factory := NewBlockerFactory(logger)
	var available []string

	for _, btype := range []BlockerType{BlockerTypeNftables, BlockerTypeIptables, BlockerTypeFirewalld, BlockerTypeUfw} {
		blocker, err := factory.Create(btype, "")
		if err == nil && blocker.IsAvailable() {
			available = append(available, blocker.Name())
		}
	}

	return available
}
