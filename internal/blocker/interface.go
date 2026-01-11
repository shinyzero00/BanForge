package blocker

type BlockerEngine interface {
	Ban(ip string) error
	Unban(ip string) error
}
