package blocker

type BlockerEngine interface {
	Ban(ip string) error
	Unban(ip string) error
	IsBanned(ip string) (bool, error)
	Flush() error
}
