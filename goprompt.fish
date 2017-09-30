function fish_prompt --description 'Write out the prompt'
set LAST_COMMAND_RC $status
goprompt -fish $GOPROMPT_OPTIONS
end
