# eatwell
Automation for the EatWell business

steps to get it running

sudo apt update -y && sudo apt upgrade -y

run through this -> https://golang.org/doc/install
export PATH=$PATH:/usr/local/go/bin
go version //to test
go get -u google.golang.org/api/docs/v1
go get -u golang.org/x/oauth2/google

go get golang.org/x/net/context
go get golang.org/x/oauth2
go get golang.org/x/oauth2/google
go get google.golang.org/api/sheets/v4

go build geoCode.go

sudo crontab -e
30 12 * * 5 cd /home/ubuntu/eatwell && ./geoCode > /home/ubuntu/eatwell/error.log
