package blocker

// BlockerEngine defines the interface for all firewall implementations
type BlockerEngine interface {
	// Core operations
	Ban(ip string) error
	Unban(ip string) error

	// Lifecycle management
	Setup() error
	Close() error

	// Query operations
	List() ([]string, error)

	// Metadata
	Name() string
	IsAvailable() bool
}
