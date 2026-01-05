package config

const Base_config = `# This is a TOML config file for BanForge it's a simple config file
# https://github.com/d3m0k1d/BanForge

# Firewall settings block
[firewall]
name = "iptables" # Name one of the support firewall(iptables, nftables, firewalld, ufw)
ban_time = 1200

[Service]
name = "nginx"
log_path = "/var/log/nginx/access.log"
enabled = true
`
