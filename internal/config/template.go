package config

const Base_config = `
# This is a TOML config file for BanForge
# [https://github.com/d3m0k1d/BanForge](https://github.com/d3m0k1d/BanForge)

[firewall]
name = ""
config = "/etc/nftables.conf"
ban_time = 1200

[[service]]
name = "nginx"
log_path = "/var/log/nginx/access.log"
enabled = true

[[service]]
name = "nginx"
log_path = "/var/log/nginx/access.log"
enabled = false

`

// TODO: fix types for use 1 or any services"
