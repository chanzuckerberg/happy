function happyvm_set_locked_version {
    eval $(hvm env)
}

add-zsh-hook chpwd happyvm_set_locked_version
