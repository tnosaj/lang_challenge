# Store
## About
This repository contains a very simple application to create and retrieve orders.

To start the app use `docker-compose` up Once started you can:

* Create a new order with `curl localhost:8080/create`. This returns the order id
* Fetch an order with `curl localhost:8080/order/YOUR_ORDER_ID`

Some facts about the application:
* It uses redis to store the orders as a json serialization.
* It is written in Golang.
* It is docker-compose to start the full stack (web server & redis)
* Follows a very weak hexagonal & DDD approach.

What this app doesnt have:
* Any tests
* Any logging functionality
* Any monitoring/metrics

## Requirements 
We would like to implement whatever you find necessary to make this code 
production ready from and SRE point of view.

## Other considerations
Ideally you would spend no more than 4 hours on this assignment
(although you can use the time you need).
There are many areas that could be covered and we would like you to
prioritize the features that you believe are more important given
the limited time budget available.

The exercise is open-ended on purpose in order to give you freedom to
experiment and show your strengths, there is not a valid answer we will
match your solution against in order to evaluate it.
