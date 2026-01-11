package blocker

import (
	"github.com/d3m0k1d/BanForge/internal/logger"
	"os/exec"
)

type Ufw struct {
	logger *logger.Logger
}

func NewUfw(logger *logger.Logger) *Ufw {
	return &Ufw{
		logger: logger,
	}
}

func (ufw *Ufw) Ban(ip string) error {
	cmd := exec.Command("sudo", "ufw", "--force", "deny", "from", ip)
	ufw.logger.Info("Banning " + ip)
	return cmd.Run()
}

func (ufw *Ufw) Unban(ip string) error {
	cmd := exec.Command("sudo", "ufw", "--force", "delete", "deny", "from", ip)
	ufw.logger.Info("Unbanning " + ip)
	return cmd.Run()
}
