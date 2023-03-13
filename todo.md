# Initial thoughts:

## Metrics
* redis latency metrics
* add counters for `err != nil` with tags for locations in code
* make metrics an interface (swap out metric in 1 place)

## Logging
* catch `err != nil` and log
* make loging an interface (swap out logging in 1 place)


## Code
* set contexts (e.g. timeouts)
* Get code should sanitize inputs
* max concurrency (429)?
* redis connection timeout
* redis connection pooling
* stop accepting on redis errors OR fail faster (background redis checks)

## Tests
* add rudimentary unit tests

## Deployment
* replace docker-compose

## DR
* redis HA
* redis backups

## Automation
* Add automatic build pipeline (github actions/gitlab ci)

# Outside the scope
* observability platform
* container runtime
