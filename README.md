# money-assist
## How to Build
* Install go v1.18+ (Please refer to golang official website)
* git clone git@gitlab.eng.vmware.com:hez/money-assist.git
* cd money-assist
* export GOPROXY=https://goproxy.cn # when you are in China
* go mod download  
* go mod vendor
* make
* ./money-assist
