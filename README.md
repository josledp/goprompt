A prompt generator writen in golang. It is heavily inspired on bash-git-prompt
(https://github.com/magicmonty/bash-git-prompt), but it adds information about
Python environment, AWS*, ...

*Aws information is based on a custom tool that exports on AWS_ROLE and
AWS_SESSION_EXPIRATION information about the current assumed role.

## Usage
* You need a valid go installation and $GOPATH/bin on your path
* You need libgit2 >= 0.25. if its not 0.26 you have to change the git2go
 library on prompt/helpers.go 
* go get github.com/josledp/goprompt
* go install github.com/josledp/goprompt
* link goprompt.sh to your home (or any other directory you may want)
* add to your .bashrc:
 source ~/goprompt.sh #Or the path you copied goprompt.sh on

## Customization
* set GOPROMPT_OPTIONS in your .bashrc with your favourites goprompt options.
  You can change it dinamically (so you can play with it until you find the
  right options for you in the console before setting it up on .bashrc)

## Plugins
* aws: uses AWS_ROLE + AWS_SESSION_EXPIRATION
* git: shows information on branch/commits diff with upstream/current workdir
  status...
* golang: shows information of the runtime golang version
* lastcommand: shows the last $?
* path: shows the current path
* python: shows current virtualenvironment if any
* user: show user@hostname (or just hostname if your are root)
* userchar: $ or # (normal user vs root)

## Known issues
* git information is not refreshed automatically (you need to run git fetch manually)
* Missing tests after plugin refactor
* The code needs a couple of iterations more to be proud of

## Todo
* Implement more plugin options
* git plugin:
- implement another styles (oh-my-zsh for example)
* Implement more plugins
- ruby? (RVM/rbenv/bundle ¿?)
- golang

