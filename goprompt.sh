setGoPrompt() {
  export LAST_COMMAND_RC=$?
  PS1=`goprompt`
}
export PROMPT_COMMAND='setGoPrompt'
