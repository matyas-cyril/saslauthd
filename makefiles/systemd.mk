define fSystemD
[Unit]
Description=Serveur d'authentification SASL 
StartLimitBurst=5
StartLimitIntervalSec=60

[Service]
Type=simple
#User=mail
#Group=mail

WorkingDirectory=${REP_INSTALL}/
ExecStart=${REP_INSTALL}/${NAME}

#WatchdogSec=10
#Restart=on-failure
#RestartSec=10

#TimeoutStopSec=10

[Install]
WantedBy=multi-user.target
endef
