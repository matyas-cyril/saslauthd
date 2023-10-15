
VERSION=0.9.0
NAME=saslauthd
ARCH=amd64

REP_PLUGINS=plugins
REP_INSTALL=/opt/saslauthd

REP_BUILD=BUILD
REP_DEB=DEB

BUILD_TIME=$(shell date +%s)

export fHelp:=$(fHelp)
export fConfFiles:=$(fConfFiles)
export fControl:=$(fControl)
export fPreInst:=${fPreInst}
export fPostInst=$(fPostInst)
export fPreRm:=${fPreRm}
export fPostRm:=${fPostRm}
export fSaslAuthdConf=$(fSaslAuthdConf)
export fSystemD=$(fSystemD)

help:
	printf "$${fHelp}"

.clean_build:
	if [ -d "${REP_BUILD}/" ]; then rm -rf ${REP_BUILD}/* && rmdir ${REP_BUILD}; fi

.clean_deb:
	if [ -d "${REP_DEB}/" ]; then rm -rf ${REP_DEB}/* && rmdir ${REP_DEB}; fi

.clean_plugins:
	if [ -d "${REP_PLUGINS}/" ]; then rm -rf ${REP_PLUGINS}/*; else mkdir ${REP_PLUGINS}; fi

clean: .clean_build .clean_deb .clean_plugins

plugins: .clean_plugins
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_PLUGINS}/random.sasl src_plugin/random/random.go && \
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_PLUGINS}/jwt.sasl src_plugin/jwt/*.go && \
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_PLUGINS}/lemon.sasl src_plugin/lemon/lemon.go && \
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_PLUGINS}/ldap.sasl src_plugin/ldap/ldap.go

man:
	rm -f ${REP_BUILD}/saslauthd.1*
	mkdir -p ${REP_BUILD} && \
	pandoc README.md -s -f markdown -t man --atx-headers -o ${REP_BUILD}/saslauthd.1 && gzip ${REP_BUILD}/saslauthd.1

build: .clean_build man

	mkdir -p ${REP_BUILD}/plugin && \
	go build -ldflags="-s -w \
	    -X github.com/matyas-cyril/saslauthd.VERSION=${VERSION} \
		-X github.com/matyas-cyril/saslauthd.BUILD_TIME=${BUILD_TIME} \
		-X github.com/matyas-cyril/saslauthd.APP_PATH=${REP_INSTALL} \
	   " \
	   -o ${REP_BUILD}/saslauthd main/main.go && \
	\
	go build -ldflags="-s -w" -buildmode=plugin -o BUILD/${REP_PLUGINS}/random.sasl src_plugin/random/random.go && \
	go build -ldflags="-s -w" -buildmode=plugin -o BUILD/${REP_PLUGINS}/jwt.sasl src_plugin/jwt/*.go && \
	go build -ldflags="-s -w" -buildmode=plugin -o BUILD/${REP_PLUGINS}/lemon.sasl src_plugin/lemon/lemon.go && \
	go build -ldflags="-s -w" -buildmode=plugin -o BUILD/${REP_PLUGINS}/ldap.sasl src_plugin/ldap/ldap.go

deb: .clean_deb build

	mkdir -p ${REP_DEB}${REP_INSTALL}/${REP_PLUGINS} \
	 		 ${REP_DEB}/lib/systemd/system \
             ${REP_DEB}/DEBIAN \
			 ${REP_DEB}/usr/share/man/man1/ 

	cp -f ${REP_BUILD}/saslauthd ${REP_DEB}${REP_INSTALL}/saslauthd
	cp -f ${REP_BUILD}/${REP_PLUGINS}/*.sasl ${REP_DEB}${REP_INSTALL}/${REP_PLUGINS}
	cp -f ${REP_BUILD}/saslauthd.1.gz ${REP_DEB}/usr/share/man/man1/saslauthd.1.gz

	printf "$${fSaslAuthdConf}" > ${REP_DEB}${REP_INSTALL}/saslauthd.toml
	printf "$${fSystemD}" > ${REP_DEB}/lib/systemd/system/saslauthd.service

	printf "$${fConfFiles}" > ${REP_DEB}/DEBIAN/conffiles && /bin/chmod 0755 ${REP_DEB}/DEBIAN/conffiles
	printf "$${fPreInst}" > ${REP_DEB}/DEBIAN/preinst && /bin/chmod 0755 ${REP_DEB}/DEBIAN/preinst
	printf "$${fPostInst}" > ${REP_DEB}/DEBIAN/postinst && /bin/chmod 0755 ${REP_DEB}/DEBIAN/postinst
	printf "$${fPreRm}" > ${REP_DEB}/DEBIAN/prerm && /bin/chmod 0755 ${REP_DEB}/DEBIAN/prerm
	printf "$${fPostRm}" > ${REP_DEB}/DEBIAN/postrm && /bin/chmod 0755 ${REP_DEB}/DEBIAN/postrm
	printf "$${fControl}" > ${REP_DEB}/DEBIAN/control && /bin/chmod 0755 ${REP_DEB}/DEBIAN/control
	
	sudo dpkg-deb -Zgzip --build ${REP_DEB}/ ${REP_DEB}/${NAME}_${VERSION}_${ARCH}.deb && \
	     sudo chown 1000:1000 ${REP_DEB}/${NAME}_${VERSION}_${ARCH}.deb

	if [ -d "${REP_DEB}/" ]; then rm -rf ${REP_DEB}/DEBIAN ${REP_DEB}/usr; fi


define fHelp
make [option]

option:

	build:
		Générer saslauthd, le fichier man et les plugins dans le dossier '${REP_BUILD}'

	clean:
		Supprimer les dossiers '${REP_BUILD}', '${REP_DEB}' et '${REP_PLUGINS}'

	deb:
		créer le paquet 'deb'

	man:
		générer le manuel de l'application à partir du fichier Readme.md

	plugins:
		compiler les plugins (.sasl) dans le répertoire '${REP_PLUGINS}'

	help:
		afficher la liste des commandes

endef

define fControl
Package: saslauthd
Version: ${VERSION}
Maintainer: matyas.cyril
Description: Serveur d'authentification SASL en Go
Section: base
Priority: optional
Architecture: amd64
Installed-Size: $(shell du -s ${REP_BUILD} 2> /dev/null | cut -f1)

endef

define fPreInst
#!/bin/bash
mkdir -p ${REP_INSTALL}

endef

define fPostInst
#!/bin/bash
chown root:root ${REP_INSTALL}/saslauthd.toml && /bin/chmod 0644 ${REP_INSTALL}/saslauthd.toml
chown root:root ${REP_INSTALL}/saslauthd && /bin/chmod 0550 ${REP_INSTALL}/saslauthd
chown -R root:root ${REP_INSTALL}/${REP_PLUGINS} && /bin/chmod -R 0440 ${REP_INSTALL}/${REP_PLUGINS}
# chown root:root /usr/share/man/man1/saslauthd.1.gz && /bin/chmod 0644 /usr/share/man/man1/saslauthd.1.gz

[ -x /bin/systemctl ] && /bin/systemctl daemon-reload

endef

define fPreRm
#!/bin/bash

[ -x /bin/systemctl ] && pgrep -l ${NAME} && /bin/systemctl stop ${NAME}

endef

define fPostRm
#!/bin/bash

[ -x /bin/systemctl ] && /bin/systemctl daemon-reload

endef

define fConfFiles
${REP_INSTALL}/saslauthd.toml

endef

define fSaslAuthdConf
[SERVER]
socket = "/var/run/saslauthd/mux"
client_max = 100
user = "cyrus"
log = "SYSLOG"

[CACHE]
enable = false
keyRand = true

[CACHE.LOCAL]
purge_on_start = true

[AUTH]
mech = ["NO"]

endef

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
ExecStart=${REP_INSTALL}/saslauthd

#WatchdogSec=10
#Restart=on-failure
#RestartSec=10

#TimeoutStopSec=10

[Install]
WantedBy=multi-user.target

endef
