### Description
This is demo golang application which shows visitor counter by GET request. 

### Stack
- golang
- mysql DB
- possibility to pass application config in environment variables (have precedence) or within configuration file
- events logging to the stdout and access/error logs files

### Requiremets
- docker
- docker-compose >= 1.25.5

### Run

```
$ docker-compose up -d
$
$
$ docker-compose ps
Name               Command               State                   Ports              
------------------------------------------------------------------------------------
app-demo    go run /app.go               Up           0.0.0.0:8099->8099/tcp          
mysql-demo  docker-entrypoint.sh mysqld  Up (healthy) 0.0.0.0:3306->3306/tcp, 33060/tcp
$
```

### URLs
- show visitors counter - **/**
```
$ curl localhost:8099
Hi!
You came from 172.29.29.1
You're 42 visiter.
Welcome!
```
- clean visitors counter - **/drop**
```
$ curl -L localhost:8099/drop
Hi!
You came from 172.29.29.1
You're 1 visiter.
Welcome!
```

