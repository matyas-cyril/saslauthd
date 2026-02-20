define fSaslAuthdConf
[SERVER]
self = false
socket = "/var/run/saslauthd/mux"
client_max = 100
log = "SYSLOG"

[CACHE]
enable = false
keyRand = true

[CACHE.LOCAL]
purge_on_start = true

[AUTH]
mech = ["NO"]

endef
