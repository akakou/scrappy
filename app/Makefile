build: 
	mkdir -p bin && gcc ./sign/*.c -o ./bin/sign `pkg-config ecdaa --libs --cflags`
install: ./attestation
	install ./attestation.json ~/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json