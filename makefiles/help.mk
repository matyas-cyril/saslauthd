define fHelp
make [option]

option:

	build:
		Générer ${NAME}, le fichier man et les plugins dans le dossier '${REP_BUILD}'

	clean:
		Supprimer les dossiers '${REP_BUILD}', '${REP_DEB}' et '${REP_PLUGINS}'

	deb:
		Créer le paquet 'deb'. Le fichier généré sera dans le dossier 'DEB'

	man:
		Générer le 'man' à partir du fichier 'man.md'

	plugins:
		Compiler les plugins (.sasl) dans le répertoire '${REP_PLUGINS}'

	help:
		Afficher la liste des commandes

endef
