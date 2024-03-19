
...

default stuff

...



# ----------------------------------------
# Michele Mendel

export PATH=$PATH:/usr/local/go/bin

alias editrc='vim ~/.bashrc'
alias sourceme='source ~/.bashrc'
alias ..='cd ..'


# -----------------------------------------------------------------------------
# Networking

alias ips='ip addr show'
alias myexternalip='curl ipecho.net/plain; echo'
alias ports1='netstat -anp tcp | grep LISTEN'
alias ports='sudo lsof -i -P -n'
alias routes='netstat -nr'


# ----------------------------------------
# DMT MMS - medlemsregister
alias dmt='cd /home/michele/dmtmms'
alias dmtstart='dmt && make server'
alias cli='dmt && ./bin/cli'


echo "You have been sourced"