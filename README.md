A prompt generator wrote in golang. It is heavily inspired on bash-git-prompt
(https://github.com/magicmonty/bash-git-prompt), but it works outside a git
repository and adds information about Python environment and AWS*

*Aws information is based on a custom tool that exports on AWS_ROLE and
AWS_SESSION_EXPIRATION information about the current assumed role.

== Usage ==
* You need a valid go installation and $GOPATH/bin on your path
* go get github.com/josledp/goprompt
* go install github.com/josledp/goprompt
* link goprompt.sh to your home (or any other directory you may want)
* add to your .bashrc:
 source ~/goprompt.sh #Or the path you copied goprompt.sh on

== Customization ==
* set GOPROMPT_OPTIONS in your .bashrc with your favourites goprompt options.
  You can change it dinamically (so you can play with it until you find the
  right options for you in the console before setting it up on .bashrc)

== Known issues==
* currently only Evermeet (Debian/Ubuntu) & Fedora styles supported, more to come
* git information is not refreshed automatically (you need to run git fetch manually)
* Missing a lot of tests
* on big repositories it is somewhat slow (probable we should cache something)
* Not happy yet with the code. Need to refactor some things
