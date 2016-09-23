FROM ntboes/golang-gin

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install lynx-cur

RUN go get gopkg.in/olivere/elastic.v2
