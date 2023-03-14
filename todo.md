# Initial thoughts:

No particular order

## Code
* ~set contexts (e.g. timeouts)~
* ~Get code should sanitize inputs~ (a bit hard with pure go, switch to router?)
* ~redis connection timeout~
* ~redis connection pooling~
* stop accepting on redis errors OR fail faster (background redis checks)
* max concurrency (http 429)?

### Metrics
* ~redis latency metrics~ [MR](https://github.com/tnosaj/lang_challenge/pull/2)
* ~add counters for `err != nil` with tags for locations in code~ [MR](https://github.com/tnosaj/lang_challenge/pull/2)
* make metrics an interface (swap out metric in 1 place)

### Logging
* ~catch `err != nil` and log~ [MR](https://github.com/tnosaj/lang_challenge/pull/3)
* make loging an interface (swap out logging in 1 place)

### Tests
* ~add rudimentary unit tests~ [MR](https://github.com/tnosaj/lang_challenge/pull/1)

## Deployment

### Redis
* ~maxmemory~ [MR](https://github.com/tnosaj/lang_challenge/pull/4)
* ~authentication~ [MR](https://github.com/tnosaj/lang_challenge/pull/4)

### runtime
* ~replace docker-compose~ [MR](https://github.com/tnosaj/lang_challenge/pull/4)

## Automation
* Add automatic build pipeline (github actions/gitlab ci)

# Outside the scope
* observability platform
* container runtime
* redis HA - related to container runtime
* redis backups - related to container runtime

