# goprompt

A plugable prompt generator writen in golang. It is heavily inspired on
[https://github.com/magicmonty/bash-git-prompt](bash-git-prompt), but it adds
information about Python environment, AWS*, golang... And you can add your
own plugins!

*Aws information is based on a custom tool that exports on AWS_ROLE and
AWS_SESSION_EXPIRATION information about the current assumed role.

## Usage

* You need a valid go installation and $GOPATH/bin on your path
* Git plugins needs libgit2 >= 0.25. if its not 0.26 you have to change the git2go
 library import on plugins/git.go (if you want to use 0.24 you have to comment
 out the stashes part on the Load function)
* go get github.com/josledp/goprompt
* go install github.com/josledp/goprompt
* For bash/zsh:
  * link goprompt.(bash|zsh) to your home (or any other directory you may want)
  * add to your .bashrc/.zshrc:
    source ~/goprompt.(bash|zsh) #Or the path you linked the file on
* For Fish:
  * link fish_prompt.fish in ~/.config/fish/functions (remove any other fish_prompt
    function you may have)

## Customization

* set GOPROMPT_OPTIONS in your .bashrc|fishd|.zshrc with your favourites
  goprompt options. You can change it dinamically (so you can play with it
  until you find the right options for you in the console before setting it up)

* There is also a configuration file at ~/.config/goprompt/goprompt.json when
  you may specify your customTemplate, and the different options a plugin may
  offer. Currently the only way to tune the plugin options is using this config
  file.

## Plugins

* aws: shows your current assumed role (red if expired, yellow if < 10minuts to
  expiration, blue if < 30 minutes else green)
* git: shows information on branch/commits diff with upstream/current workdir
  status.... It does a fetch if last fetch >300 seconds
* golang: shows information of the runtime golang version
* lastcommand: shows the last command return code
* path: shows the current path
* python: shows current virtualenvironment if any
* user: shows the user (if its not root)
* hostname: shows the hostname (green if regular user, red if root)
* userchar: $ or # (normal user vs root)

## Known issues
* Missing some tests after plugin refactor
* The code needs a couple of iterations more to be proud of

## Todo
* Implement plugin options in the command line
* git plugin:
  * implement another styles (oh-my-zsh for example)
* Implement more plugins:
  * ruby? (RVM/rbenv/bundle ¿?)
* Setup more predefined templates
