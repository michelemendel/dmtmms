
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
alias editservice='sudo vim /etc/systemd/system/dmtmms.service'
alias daemonreload='sudo systemctl daemon-reload'
alias dmtstart='sudo systemctl start dmtmms'
alias dmtrestart='sudo systemctl restart dmtmms'
alias dmtstop='sudo systemctl stop dmtmms'
alias dmtstatus='sudo systemctl status dmtmms'
alias dmtenable='sudo systemctl enable dmtmms'
alias listservices='sudo systemctl list-units --type=service'
alias dmttail='dmt && tail -f log/log.log'
alias cli='dmt && ./bin/cli'
#alias dmtstart='dmt && make server &'
#alias dmtstop='ps aux | grep .bin/server | grep -v grep | awk "{print \$2}" | xargs kill'

dmt

echo "You have been sourced"