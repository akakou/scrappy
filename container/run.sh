mkdir -p /home/attestation/.config/google-chrome/NativeMessagingHosts/
cp /attestation/client/app/attestation.json /home/attestation/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
chown attestation /home/attestation/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json

### build

### Chrome
chown attestation /dev/tpm0
chown attestation -R /attestation/

cd /attestation/
su attestation -c 'google-chrome --no-sandbox -user-data-dir="/home/attestation/.config/google-chrome" http://localhost:5000'


cd /attestation/server
python3 run.py

