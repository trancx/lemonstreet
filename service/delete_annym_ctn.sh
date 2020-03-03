docker container ls --all | grep -Eo "[a-z]+_[a-z]+" | xargs docker container rm 
