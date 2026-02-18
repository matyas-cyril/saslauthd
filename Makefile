
VERSION=1.0.0
NAME=go-saslauthd
ARCH=amd64

REP_PLUGINS=plugins
REP_INSTALL=/opt/go-saslauthd

DEFAULT_CONF_FILE=saslauthd.toml

REP_DEST=/tmp

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

-include makefiles/*.mk


help:
	printf "$${fHelp}"

.var_plugins:
	$(eval REP_DEST := plugins)

.var_plugins_build:
	$(eval REP_DEST := BUILD/plugins)

.rm_build: .var_plugins_build
	@if [ -d "${REP_BUILD}/" ]; then rm -rf ${REP_BUILD}/* && rmdir ${REP_BUILD}; fi

.rm_deb:
	@if [ -d "${REP_DEB}/" ]; then rm -rf ${REP_DEB}/* && rmdir ${REP_DEB}; fi

.rm_plugins: .var_plugins
	@if [ -d "${REP_DEST}/" ]; then rm -rf ${REP_DEST}/*; else mkdir ${REP_DEST}; fi

.rm_pgauth:
	@if [ -d "${REP_DEST}/" ]; then rm -rf ${REP_DEST}/pgauth.sasl; else mkdir ${REP_DEST}; fi

.rm_random: .var_plugins
	@if [ -d "${REP_DEST}/" ]; then rm -rf ${REP_DEST}/random.sasl; else mkdir ${REP_DEST}; fi

.rm_jwt: .var_plugins
	@if [ -d "${REP_DEST}/" ]; then rm -rf ${REP_DEST}/jwt.sasl; else mkdir ${REP_DEST}; fi

.rm_lemon: .var_plugins
	@if [ -d "${REP_DEST}/" ]; then rm -rf ${REP_DEST}/lemon.sasl; else mkdir ${REP_DEST}; fi

.rm_ldap: .var_plugins
	@if [ -d "${REP_DEST}/" ]; then rm -rf ${REP_DEST}/ldap.sasl; else mkdir ${REP_DEST}; fi

clean: .rm_build .rm_deb .rm_plugins

.pgauth:
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_DEST}/pgauth.sasl src_plugin/pgauth/pgAuth.go src_plugin/pgauth/define.go

pgauth: .var_plugins .rm_pgauth .pgauth
	
.random:
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_DEST}/random.sasl src_plugin/random/random.go

random: .var_plugins .rm_random .random

.jwt:
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_DEST}/jwt.sasl src_plugin/jwt/*.go

jwt: .var_plugins .rm_jwt .jwt

.lemon:
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_DEST}/lemon.sasl src_plugin/lemon/lemon.go

lemon: .var_plugins .rm_lemon .lemon

.ldap:
	go build -ldflags="-s -w" -buildmode=plugin -o ${REP_DEST}/ldap.sasl src_plugin/ldap/ldap.go

ldap: .var_plugins .rm_ldap .ldap

plugins: .var_plugins .rm_plugins pgauth random jwt ldap

.go_clean:
	go clean

.build_sasl: 
	mkdir -p ${REP_BUILD}/plugins && \
	go install && \
	go build -trimpath -ldflags="-s -w \
	    -X github.com/matyas-cyril/saslauthd.VERSION=${VERSION} \
		-X github.com/matyas-cyril/saslauthd.BUILD_TIME=${BUILD_TIME} \
		-X github.com/matyas-cyril/saslauthd.APP_NAME=${NAME} \
		-X github.com/matyas-cyril/saslauthd.APP_CONF=${REP_INSTALL}/${DEFAULT_CONF_FILE} \
	   " \
	   -o ${REP_BUILD}/${NAME} main/main.go 

man:
	rm -f ${REP_BUILD}/${NAME}.1*
	mkdir -p ${REP_BUILD} && \
	pandoc man.md -s -f markdown -t man -o ${REP_BUILD}/${NAME}.1 && gzip ${REP_BUILD}/${NAME}.1

build: .rm_build man .go_clean .build_sasl .var_plugins_build .jwt .random .ldap .lemon .pgauth

deb: .rm_deb build

	mkdir -p ${REP_DEB}${REP_INSTALL}/${REP_PLUGINS} \
	 		 ${REP_DEB}/lib/systemd/system \
             ${REP_DEB}/DEBIAN \
			 ${REP_DEB}/usr/share/man/man1

	cp -f ${REP_BUILD}/${NAME} ${REP_DEB}${REP_INSTALL}/${NAME}
	cp -f ${REP_BUILD}/${REP_PLUGINS}/*.sasl ${REP_DEB}${REP_INSTALL}/${REP_PLUGINS}
	cp -f ${REP_BUILD}/${NAME}.1.gz ${REP_DEB}/usr/share/man/man1/${NAME}.1.gz

	printf "$${fSaslAuthdConf}" > ${REP_DEB}${REP_INSTALL}/saslauthd.toml
	printf "$${fSystemD}" > ${REP_DEB}/lib/systemd/system/${NAME}.service

	printf "$${fConfFiles}" > ${REP_DEB}/DEBIAN/conffiles && /bin/chmod 0755 ${REP_DEB}/DEBIAN/conffiles
	printf "$${fPreInst}" > ${REP_DEB}/DEBIAN/preinst && /bin/chmod 0755 ${REP_DEB}/DEBIAN/preinst
	printf "$${fPostInst}" > ${REP_DEB}/DEBIAN/postinst && /bin/chmod 0755 ${REP_DEB}/DEBIAN/postinst
	printf "$${fPreRm}" > ${REP_DEB}/DEBIAN/prerm && /bin/chmod 0755 ${REP_DEB}/DEBIAN/prerm
	printf "$${fPostRm}" > ${REP_DEB}/DEBIAN/postrm && /bin/chmod 0755 ${REP_DEB}/DEBIAN/postrm
	printf "$${fControl}" > ${REP_DEB}/DEBIAN/control && /bin/chmod 0755 ${REP_DEB}/DEBIAN/control
	
	fakeroot dpkg-deb -Zgzip --build ${REP_DEB}/ ${REP_DEB}/${NAME}_${VERSION}_${ARCH}.deb && \
	     fakeroot chown 1000:1000 ${REP_DEB}/${NAME}_${VERSION}_${ARCH}.deb

	@if [ -d "${REP_DEB}/" ]; then rm -rf ${REP_DEB}/DEBIAN ${REP_DEB}/usr; fi
