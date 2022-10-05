# rm -rf /home/attestation/.config
mkdir -p /home/attestation/.config

mkdir -p /home/attestation/.config/google-chrome/NativeMessagingHosts/
cp /attestation/client/app/attestation.json /home/attestation/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
chown attestation -R /home/attestation/.config/

### build
export PATH=$PATH:/usr/local/go/bin
export GOFLAGS='-buildvcs=false'

cd /attestation/client/app
go build 

### init
cd /attestation/init
go build 
./m

### Chrome
chown attestation /dev/tpm0
chown attestation -R /attestation/

cd /attestation/
su attestation -c 'google-chrome --no-sandbox -user-data-dir="/home/attestation/.config/google-chrome" http://localhost:8080  &'

cd /attestation/verifier/
go build && go run main.go


