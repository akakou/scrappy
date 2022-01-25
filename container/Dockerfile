FROM --platform=linux/amd64 ubuntu

RUN apt-get update && apt-get install -y wget tzdata 
RUN apt-get install -y gnupg

RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list
RUN apt-get update && apt-get install -y google-chrome-stable

RUN apt-get install -y python3 python3-pip
#RUN pip3 install selenium chromedriver-binary sqlalchemy
RUN apt-get install gcc

RUN useradd -m -u 1000 attestation
RUN mkdir -p /home/attestation
RUN chown attestation /home/attestation -R
USER attestation


