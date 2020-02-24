BUG

il y a un bug dans fyne, qui utilise vendor/mobile/x
or utiliser ce vendor cree 2 references differentes au context android (l'autre avec golang/...) et a cause de cela,
on ne peut pas utiliser le package asset pour acceder a des resources dans le apk:
cf. https://github.com/golang/go/issues/26445
la solution est de ne pas vendoriser mobile/x/ mais est-ce que alors mes fix marcheront toujours ?

TODO:

- on lui donne dans un fichier/res l'URL d'une outbox: comment y acceder depuis un res dans APK ?
- il la charge pour recuperer le "last" URL
- il download an arriere-plan a l'init tous les posts, et en garde une copie locale
- il faut ensuite gerer les downloads des nouveaux posts seulement, toujours en partant du "last"
- il faut afficher la date et l'auteur pour chaque post
- lorsque l'appli est relancee, il faut afficher a partir du dernier post lu

- faire entrer un code au user pour qu'il puisse ensuite se connecter grace au code et pouvoir pusher un message sur l'AP server
