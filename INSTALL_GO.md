# install

## Mac

    $ brew install go
    $ brew install direnv
    
# move to workspace

    $ cat ~/.bash_profile
    export PATH="/usr/local/opt/sqlite/bin:$PATH"
    export EDITOR=vim
    eval "$(direnv hook bash)"
    $ cd /path/to/dir
    $ direnv edit .
