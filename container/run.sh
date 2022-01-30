### Chrome
mkdir -p /home/attestation/.config/google-chrome/NativeMessagingHosts
cp /home/attestation/attestation/app/attestation.json /home/attestation/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
chmod +x /home/attestation/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
google-chrome --no-sandbox -user-data-dir="/home/attestation/.config/google-chrome"


