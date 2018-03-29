# goprompt

A plugable prompt generator written in golang. It is heavily inspired on
[https://github.com/magicmonty/bash-git-prompt](bash-git-prompt), but it adds
information about Python environment, AWS*, golang... And you can add your
own plugins!

*Aws information is based on a custom tool that exports on AWS_ROLE and
AWS_SESSION_EXPIRATION information about the current assumed role.

## Usage

* You need a valid go installation and $GOPATH/bin on your path
* libgit2 version 0.26 (Git plugin needs git2go, which is has the go bindings for libgit2).
  Currently using git2go.v26 (bindings agains libgit2 0.26) if you have another
  version you will have to change the import to be able to build goprompt.
  * (on MAC you can install libgit2 with brew. you will need pkg-config if its not already installed)
* go get github.com/josledp/goprompt
* go install github.com/josledp/goprompt
* For bash/zsh:
  * link goprompt.(bash|zsh) (on the repository root) to your home (or any other directory you may want)
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
  Example:
    ```{"options":{"path.fullpath":2}}```
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
