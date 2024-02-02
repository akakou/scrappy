# SeCure Rate Assuring Protocol with PrivacY
  
## Dependencies

Scrappy is support Linux environments, and the work checked on  NixOS 23.05pre4.
Also, Scrappy needs the dependecies as follows:

- docker and docker-compose
- tpm-tools (https://github.com/tpm2-software/tpm2-tools)


## How to use
### 1. Reset TPM

**WARNING: TO DEPLOY CORRECTLY, WE RECOMMEND RESETTING YOUR TPM FIRST. BUT IT MAY REMOVE YOUR IMPORTANT KEYS.**

```sh
sudo tpm2_clear
```
  
### 2. Set up X configuration
  
```sh
xhost + 
```

### 3. Download libraries

```
cd thirdparty
make
cd ..
```
  
### 4. Run docker-compose
  
The run follows the command, then the Chrome browser opens.
  
```
docker-compose --profile default up
```
  
### 5. Configure chrome extension
  
1. Open new tabs and jump to "chrome://extensions/" on the browser.
2. Turn on the toggle of "Developer mode".
3. Push the button labeled and open "/scrappy/browser/extension/" to install the extension.
  
## Demo
  
"http://server:8081/" is a demo that protects heavy endpoints using the Scrappy.
  
## Benchmark
  
You can run the "go test" with some commands to measure benchmarks.
  
### Cryptographic process of Sign TPM
  
```sh
cd /scrappy/ecdaa
/usr/local/go/bin/go test ./bench -benchmem -run=^$ -bench SignTPM -benchtime 20x
```
  
### Sign Log
  
```sh
cd /scrappy/core
/usr/local/go/bin/go test ./tests -benchmem -run=^$ -bench BenchmarkSignLog -benchtime 20x
```
  
### Cryptographic process of Verification
  
```sh
cd /scrappy/ecdaa
go test ./bench -benchmem -run=^$ -bench BenchmarkVerify -benchtime 20x
```
  
### Verify Log
  
```sh
cd /scrappy/core
go test ./tests -benchmem -run=^$ -bench BenchmarkVerifyLog -benchtime 20x
```
