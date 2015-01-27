GoPics
=======
![A screen of the user timeline in GoPics](https://raw.githubusercontent.com/lucachr/gopics/master/screen.png)
GoPics is a simple image sharing board built with Go, Redis, UIkit, and 
good intentions.  

The project goal is to built a simple Go web app that is bigger than the 
classics wiki or WebSocket chat examples but still small enough to be useful
for learning purpose.  

This is the first version of the app, it features only login, forms validation, 
and images upload capabilities. Real time stuff with WebSocket coming soon!

Installation
-------------

```shell
    go get github.com/lucachr/gopics
```
Usage
------

[Redis](http://redis.io/) have to be installed and a Redis instance needs to be 
up and running.  
After installing GoPics, start the server with

```shell
   $ gopics
```
and go to [localhost:8080](http://localhost:8080).

License
--------

GoPics is released under [the MIT License](http://opensource.org/licenses/MIT).
