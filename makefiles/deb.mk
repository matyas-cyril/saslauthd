define fControl
Package: ${NAME}
Version: ${VERSION}
Maintainer: matyas-cyril
Description: Serveur d'authentification SASL en Go
Section: base
Priority: optional
Architecture: amd64
Installed-Size: $(shell du -s ${REP_BUILD} 2> /dev/null | cut -f1)
Depends: man

endef

define fPreInst
#!/bin/bash
mkdir -p ${REP_INSTALL}

endef

define fPostInst
#!/bin/bash
chown mail:mail ${REP_INSTALL}/saslauthd.toml && /bin/chmod 0640 ${REP_INSTALL}/saslauthd.toml
chown mail:mail ${REP_INSTALL}/${NAME} && /bin/chmod 0550 ${REP_INSTALL}/${NAME}
chown -R mail:mail ${REP_INSTALL}/${REP_PLUGINS} && /bin/chmod -R 0440 ${REP_INSTALL}/${REP_PLUGINS}
chown root:root /usr/share/man/man1/${NAME}.1.gz && /bin/chmod 0644 /usr/share/man/man1/${NAME}.1.gz

[ -x /bin/systemctl ] && /bin/systemctl daemon-reload || exit 0

endef

define fPreRm
#!/bin/bash

[ -x /bin/systemctl ] && pgrep -l ${NAME} && /bin/systemctl stop ${NAME} || exit 0

endef

define fPostRm
#!/bin/bash

[ -x /bin/systemctl ] && /bin/systemctl daemon-reload || exit 0

endef

define fConfFiles
${REP_INSTALL}/saslauthd.toml

endef
