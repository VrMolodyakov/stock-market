# Stock service
![image](https://user-images.githubusercontent.com/99216816/211168525-c1df49d6-2350-42c7-bd5f-7a214293d047.png)



### Table of contents
* [General info](#Generalinfo)
* [Technologies](#Technologies)
* [Build](#Build)
* [Start](#Start)
* [Test](#Test)
* [Metrics](#Metrics)


## General info
This is a project that allows you to view information about the stock price of the companies of the current or last closed trades.
![image](https://user-images.githubusercontent.com/99216816/211168603-adb1ac6f-4431-4f99-9544-c849c2ce03fb.png)



## Technologies
Project is created with:
* Golang version: 1.18
* React version: 18.2.0
* Grafana version: 6.1.6
* Prometheus version: 6.1.6
* Redis version: 6.2
* Postgres version: 12.0
* Docker



## Build
To build this project, you need to run the following commands:

```
https://github.com/VrMolodyakov/stock-market.git
make build
```


## Start
To run this project, you need to run the following command:

```
make start
```
and then visit:
```
http://localhost:3001/
```



## Test
To run tests:

```
make test
```

| package | Coverage |
| ------ | ------ |
| tokenStorage   |  94.1% of statements |
| userStorage | 88.9% of statements |
| auth | 90.5% of statements |
| middleware   | 73.7% of statements |
| stock handler  | 76.6% of statements |
| service  | 100.0% of statements |



## Metrics
To setup Grafana monitoring visit ```http://localhost:3000/```.
Example:
![image](https://user-images.githubusercontent.com/99216816/211171034-3bc777f5-7941-42fc-a866-1b8178fcc43c.png)


Prometheus is available at:

```
http://localhost:9090/
```


#TL;DR
The project was made for the purpose of learning , work practice with various tools such as React,Prometheus, Grafana, Redis and so on.
