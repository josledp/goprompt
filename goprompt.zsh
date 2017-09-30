precmd() {
  export LAST_COMMAND_RC=$?
  PS1=`goprompt $GOPROMPT_OPTIONS`
}
