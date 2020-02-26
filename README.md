TODO:

Fonctionnalités:

1- client OLKi/Masto/Writeas: afficher/naviguer/rédiger commentaires
2- client TALN: agenda, news, abstracts, commentaires, likes
3- authentification par un mot de passe prédéfini (simple, supporté seulement par mon Writeas)

Fcts techniques:

- Naviguer dans une outbox (est-ce suffisant, ou faut-il gérer les timelines ?)

Naviguer dans outbox:

- charger la suite de la liste à la demande
- sauver tout ce qui est chargé dans fichier pour mode offline

Ergonomie:

- Comment naviguer entre les outbox ? bookmarks, liens, search

GUI:

X séparer les posts
- réduire/ajuster font size
- afficher les premières lignes d'un post en bas: pas necessaire, car on affiche au minimum un post, et si on clic dessus, on le fera apparaitre en entier; mais il faudra ajouter un scrollbar !!
- ajouter icône pour les URL
- entrer mot de passe pour poster des commentaires dans le fedi via mon writeas

Authentification:

Le but est de rendre l'appli la plus simple possible, donc pas de registration/login:
je crée un serveur relai Writeas avec une liste de mots de passe prédéfinis, un par user,
et je distribute les mots de passe par email.
Ca ne scale pas pour distributer l'appli à grande échelle, mais ça suffit pour une distribution ciblée.

BUG

il y a un bug dans fyne, qui utilise vendor/mobile/x
or utiliser ce vendor cree 2 references differentes au context android (l'autre avec golang/...) et a cause de cela,
on ne peut pas utiliser le package asset pour acceder a des resources dans le apk:
cf. https://github.com/golang/go/issues/26445
la solution est de ne pas vendoriser mobile/x/ mais est-ce que alors mes fix marcheront toujours ?


