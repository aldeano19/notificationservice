# Notification Service

## Overview
A notification service to send email and text messages

## Technologies used
* Golang
* Docker
* Elasticsearch


RUN with Docker (linux)
============
* cd into this repo's local version
$ sudo docker-compose build
$ sudo docker-compose up


REQUIREMENTS for production(without docker):
============
install java:
    * sudo add-apt-repository -y ppa:webupd8team/java
    * sudo apt-get update
    * sudo apt-get -y install oracle-java8-installer
elasticsearch:     
    * wget https://download.elastic.co/elasticsearch/elasticsearch/elasticsearch-1.7.2.deb
    *sudo dpkg -i elasticsearch-1.7.2.deb

      -NOTE: This results in Elasticsearch being installed in /usr/share/elasticsearch/ with its configuration files placed in /etc/elasticsearch and its init script added in /etc/init.d/elasticsearch
      
      -(Optional) To make sure Elasticsearch starts and stops automatically with the Droplet, add its init script to the default runlevels with the command:    $ sudo update-rc.d elasticsearch defaults

    * add lines to /etc/elasticsearch/elasticsearch.yaml
        http.cors.enabled: true
        http.cors.allow-origin: "*"

    * Run with: $ sudo /usr/share/elasticsearch/bin/elasticsearch & >/dev/nul

    * Copy line to /etc/hosts:
        127.0.0.1 elasticsearch


EXAMPLES
==========

SEND REQUEST
http://172.17.0.3:3000/getbrokenlinks?domain=71lbs.com

resp: {"BrokenLinks":[],"Domain":"yahooss.com","Status":"running","StartTime":"19:11:32","LastUpdateTime":"19:11:32","Error":"None","TicketId":"AVYUBRZZovZXueLikNuJ","TotalInternalLinks":0}

-------------------------------------------------------------------------------

GET REPORT
http://172.17.0.3:3000/brokenlinksreport?ticketid=AVYT5g6Li1lVGFrTGNHT

resp: {"BrokenLinks":[],"Domain":"yahoo.com","Status":"success","StartTime":"18:50:53","LastUpdateTime":"18:52:02","Error":"None","TicketId":"AVYT8i3sovZXueLikNuD","TotalInternalLinks":58}
