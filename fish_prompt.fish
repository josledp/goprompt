function fish_prompt --description 'Write out the prompt'
set LAST_COMMAND_RC $status
goprompt $GOPROMPT_OPTIONS
end
