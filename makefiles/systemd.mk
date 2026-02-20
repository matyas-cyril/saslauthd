define fSystemD
[Unit]
Description=Golang SASL Authentication Server 
StartLimitBurst=3
StartLimitIntervalSec=60

[Service]
Type=simple

User=mail
Group=mail

WorkingDirectory=${REP_INSTALL}/

RuntimeDirectory=saslauthd
RuntimeDirectoryMode=0760

ExecStart=${REP_INSTALL}/${NAME}

WatchdogSec=10
Restart=on-failure
RestartSec=10

TimeoutStopSec=10

ProtectSystem=strict
ProtectHome=yes
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target

endef
