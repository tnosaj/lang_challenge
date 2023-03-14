# Initial thoughts:

No particular order

## Code
* ~set contexts (e.g. timeouts)~
* ~Get code should sanitize inputs~ (a bit hard with pure go, switch to a non standard router?)
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
* ~Add automatic build pipeline (github actions/gitlab ci)~ [MR](https://github.com/tnosaj/lang_challenge/pull/5) (not a huge fan of github actions yet...)

# Outside the scope
* observability platform - e.g. https://github.com/prometheus-operator/prometheus-operator ... pod/svc monitors already exist
* container runtime - full fledged k8s? I tested on minikube
* redis HA - related to container runtime - e.g. https://github.com/ot-container-kit/redis-operator ... not a lot of information about redis operators honestly
* redis backups - related to container runtime - disk snapshots, replica snapshots, etc.

