setGoPrompt() {
  export LAST_COMMAND_RC=$?
  PS1=`goprompt $GOPROMPT_OPTIONS`
}
export PROMPT_COMMAND='setGoPrompt'
