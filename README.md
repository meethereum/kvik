# kvik
a distributed key value store in golang


## start the db server : 

1. make sure you have go installed

2. to start the server :  

```shell
bash launch.sh
```

3. set key value pair using the API : 
```shell
curl 'http://localhost:8090/set?key=a&value=c'
```

4. get the value using the API :
```shell
curl 'http://localhost:8090/get?key=a'
```

- next todo:

1.  [X] add gitignore
2.  [ ] complete sharding logic : use static sharding and add redirection functionality
3.  [X] make setup script
