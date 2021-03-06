#/bin/bash

origin=$('pwd')
mkdir /tmp/go
cd /tmp/go

arch=`uname -m`

if [[ "$arch" == 'x86_64' ]]
then
    wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
else
    wget https://storage.googleapis.com/golang/go1.4.2.linux-386.tar.gz
    sudo tar -C /usr/local -xzf go1.4.2.linux-386.tar.gz
fi

cd -
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile

go run hello.go

cd $origin
go install github.com/jaodsilv/cos
