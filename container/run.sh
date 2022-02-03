cp /home/attestation/client/attestation.json ~/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json

### build
cd /attestation/

### Chrome
cd /attestation/client/app
make install
google-chrome --no-sandbox -user-data-dir="/home/attestation/.config/google-chrome" http://localhost:5000

cd /attestation/server
python3 run.py


