A prompt generator wrote in golang. It is heavily inspired on bash-git-prompt
(https://github.com/magicmonty/bash-git-prompt) but adding information about
Python environment and AWS*

*Aws information is based on a custom tool that exports on AWS_ROLE and
AWS_SESSION_EXPIRATION information about the current assumed role.

== Usage ==
* You need a valid go installation and $GOPATH/bin on your path
* go get github.com/josledp/goprompt
* go install github.com/josledp/goprompt
* copy goprompt.sh to your home (or any other directory you may want)
* add to your .bashrc:
 source ~/goprompt.sh #Or the path you copied goprompt.sh on

== Customization ==
* set GOPROMPT_OPTIONS in your .bashrc with your favourites goprompt options (none yet :))

== Known issues==
* currently only Evermeet (Debian/Ubuntu) style supported, more to come
* git information is not refreshed automatically (you need to run git fetch manually)
* Missing some tests
* on big repositories it is somewhat slow (probable we should cache something)

