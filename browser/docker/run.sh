# rm -rf /home/scrappy/.config
mkdir -p /home/scrappy/.config

mkdir -p /home/scrappy/.config/google-chrome/NativeMessagingHosts/
cp browser/app/scrappy.json /home/scrappy/.config/google-chrome/NativeMessagingHosts/com.akakou.scrappy.json
chown scrappy -R /home/scrappy/.config/

### build
export PATH=$PATH:/usr/local/go/bin
export GOFLAGS='-buildvcs=false'

cd /scrappy/browser/app
go build 

### Chrome
chown scrappy /dev/tpm0
chown scrappy -R /scrappy/

cd /scrappy/
su scrappy -c 'google-chrome --no-sandbox -user-data-dir="/home/scrappy/.config/google-chrome" http://server:8081'



