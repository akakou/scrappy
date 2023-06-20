### chrome
wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add -
echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list
apt-get update && apt-get install -y google-chrome-stable

apt-get install -y python3 python3-pip
#pip3 install selenium chromedriver-binary sqlalchemy
apt-get install gcc

useradd -m -u 1000 scrappy
mkdir -p /scrappy
chown scrappy /scrappy -R
