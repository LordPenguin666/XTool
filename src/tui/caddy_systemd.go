package tui

//func (ui *UI) CaddySystemdReplacement() *UI {
//	service := "# caddy.service\n" +
//		"#\n" +
//		"# For using Caddy with a config.json file.\n" +
//		"#\n" +
//		"# Make sure the ExecStart and ExecReload commands are correct\n" +
//		"# for your installation.\n" +
//		"#\n" +
//		"# See https://caddyserver.com/docs/install for instructions.\n" +
//		"#\n" +
//		"# WARNING: This service does not use the --resume flag, so if you\n" +
//		"# use the API to make changes, they will be overwritten by the\n" +
//		"# Caddyfile next time the service is restarted. If you intend to\n" +
//		"# use Caddy's API to configure it, add the --resume flag to the\n" +
//		"# `caddy run` command or use the caddy-api.service file instead.\n" +
//		"\n" +
//		"[Unit]\n" +
//		"Description=Caddy web server\n" +
//		"Documentation=https://caddyserver.com/docs/\n" +
//		"After=network-online.target\n" +
//		"Wants=network-online.target systemd-networkd-wait-online.service\n" +
//		"StartLimitIntervalSec=14400\n" +
//		"StartLimitBurst=10\n" +
//		"\n" +
//		"[Service]\n" +
//		"Type=notify\n" +
//		"User=caddy\n" +
//		"Group=caddy\n" +
//		"Environment=XDG_DATA_HOME=/var/lib\n" +
//		"Environment=XDG_CONFIG_HOME=/etc\n" +
//		"ExecStartPre=/usr/bin/caddy validate --config /etc/caddy/config.json --adapter jsonc\n" +
//		"ExecStart=/usr/bin/caddy run --environ --config /etc/caddy/config.json --adapter jsonc\n" +
//		"ExecReload=/usr/bin/caddy reload --config /etc/caddy/config.json --adapter jsonc --force\n" +
//		"ExecStopPost=/usr/bin/rm -f /run/caddy/admin.socket\n" +
//		"LimitNOFILE=1048576\n" +
//		"LimitNPROC=512\n" +
//		"PrivateDevices=yes\n" +
//		"PrivateTmp=true\n" +
//		"ProtectSystem=full\n" +
//		"AmbientCapabilities=CAP_NET_BIND_SERVICE\n" +
//		"\n" +
//		"# Do not allow the process to be restarted in a tight loop. If the\n" +
//		"# process fails to start, something critical needs to be fixed.\n" +
//		"Restart=on-abnormal\n" +
//		"\n" +
//		"# Use graceful shutdown with a reasonable timeout\n" +
//		"TimeoutStopSec=5s"
//
//	if err := os.WriteFile("/lib/systemd/system/caddy.service", []byte(service), 644); err != nil {
//		panic(err)
//	}
//
//	return ui
//}
