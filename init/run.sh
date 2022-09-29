echo 'Remove all data in TPM!!!!!!!!'
echo 'If you are ok, press any key to continue'

read

sudo tpm2_clear

go build 

sudo ./m > config.json
