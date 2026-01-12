package blocker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type Ufw struct {
	logger *logger.Logger
}

func NewUfw(logger *logger.Logger) *Ufw {
	return &Ufw{
		logger: logger,
	}
}

// Name returns the blocker engine name
func (u *Ufw) Name() string {
	return "ufw"
}

// IsAvailable checks if ufw is available in the system
func (u *Ufw) IsAvailable() bool {
	cmd := exec.Command("which", "ufw")
	return cmd.Run() == nil
}

// Setup initializes UFW (if not already enabled)
func (u *Ufw) Setup() error {
	// Check if UFW is enabled
	cmd := exec.Command("sudo", "ufw", "status")
	output, err := cmd.CombinedOutput()

	if err != nil || !strings.Contains(string(output), "active") {
		u.logger.Warn("UFW is not active, attempting to enable...")
		cmd := exec.Command("sudo", "ufw", "--force", "enable")
		output, err := cmd.CombinedOutput()
		if err != nil {
			u.logger.Error("failed to enable UFW",
				"error", err.Error(),
				"output", string(output))
			return fmt.Errorf("failed to enable UFW: %w", err)
		}
		u.logger.Info("UFW enabled successfully")
	}

	u.logger.Info("UFW setup completed")
	return nil
}

// Ban adds an IP to the deny list
func (u *Ufw) Ban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "ufw", "--force", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		u.logger.Error("failed to ban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return fmt.Errorf("failed to ban IP %s: %w", ip, err)
	}

	u.logger.Info("IP banned", "ip", ip, "output", string(output))
	return nil
}

// Unban removes an IP from the deny list
func (u *Ufw) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "ufw", "--force", "delete", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		u.logger.Error("failed to unban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return fmt.Errorf("failed to unban IP %s: %w", ip, err)
	}

	u.logger.Info("IP unbanned", "ip", ip, "output", string(output))
	return nil
}

// List returns all currently denied IPs
func (u *Ufw) List() ([]string, error) {
	cmd := exec.Command("sudo", "ufw", "status", "numbered")
	output, err := cmd.CombinedOutput()
	if err != nil {
		u.logger.Error("failed to list UFW rules",
			"error", err.Error(),
			"output", string(output))
		return nil, fmt.Errorf("failed to list UFW rules: %w", err)
	}

	var deniedIPs []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Looking for lines with "Deny" and "From"
		if strings.Contains(line, "Deny") && strings.Contains(line, "Anywhere on") {
			// Parse UFW status output format
			parts := strings.Fields(line)
			for i, part := range parts {
				// Extract IP that comes after "from"
				if part == "from" && i+1 < len(parts) {
					ip := parts[i+1]
					if validateIP(ip) == nil {
						deniedIPs = append(deniedIPs, ip)
					}
					break
				}
			}
		}
	}

	return deniedIPs, nil
}

// Close performs cleanup operations (placeholder for future use)
func (u *Ufw) Close() error {
	// No cleanup needed for UFW
	u.logger.Info("UFW blocker closed")
	return nil
}
