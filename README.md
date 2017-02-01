<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Overview](#overview)
- [Building and Testing](#building-and-testing)
- [Configuration](#configuration)
  - [The configuration format](#the-configuration-format)
    - [Definitions](#definitions)
    - [Format overview](#format-overview)
      - [Example 1](#example-1)
      - [Example 2](#example-2)
      - [Example 3](#example-3)
    - [Descriptor list definition](#descriptor-list-definition)
    - [Rate limit definition](#rate-limit-definition)
  - [Loading Configuration](#loading-configuration)
- [Rate limit statistics](#rate-limit-statistics)
- [Debug Port](#debug-port)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Overview

The rate limit service is a Go/GRPC service designed to enable generic rate limit scenarios from different types of
applications. Applications request a rate limit decision based on a domain and a set of descriptors. The service
reads the configuration from disk via [runtime](https://github.com/lyft/goruntime), composes a cache key, and talks to the redis cache. A
decision is then returned to the caller.

# Building and Testing

* Install redis-server
* Make sure go is setup correctly and checkout rate limit service into your go path
* In order to run the integration tests using a local default redis install you will also need these environment variables set:
```
export REDIS_SOCKET_TYPE=tcp
export REDIS_URL=localhost:6379
```
* To setup for the first time (only done once):
```
make bootstrap
```
* To compile:
```
make compile
```
* To compile and run tests:
```
make tests
```
* To run the server locally using some sensible default settings you can do this (this will setup the server to read the configuration files from the path you specify):
```
USE_STATSD=false LOG_LEVEL=debug REDIS_SOCKET_TYPE=tcp REDIS_URL=localhost:6379 RUNTIME_ROOT=/home/user/src/runtime/data RUNTIME_SUBDIRECTORY=ratelimit
```

# Configuration

## The configuration format

### Definitions

* **Domain:** A domain is a container for a set of rate limits. All domains known to the rate limit service must be
globally unique. They serve as a way for different teams/projects to have rate limit configurations that don't conflict.
* **Descriptor:** A descriptor is a list of key/value pairs owned by a domain that the rate limit service uses to
select the correct rate limit to use when limiting. Descriptors are case-sensitive. Examples of descriptors are:
  * ("database", "users")
  * ("message_type", "marketing"),("to_number","2061234567")
  * ("to_cluster", "service_a")
  * ("to_cluster", "service_a"),("from_cluster", "service_b")

### Format overview

#### Example 1

Let's start with a simple example:

```yaml
domain: mongo_cps
descriptors:
  - key: database
    value: users
    rate_limit:
      unit: second
      requests_per_unit: 500

  - key: database
    value: default
    rate_limit:
      unit: second
      requests_per_unit: 500
```

The rate limit configuration file format is YAML (mainly so that comments are supported). In the configuration above
the domain is "mongo_cps" and we setup 2 different rate limits in the top level descriptor list. Each of the limits
have the same key ("database"). They have a different value ("users", and "default"), and each of them setup a 500
request per second rate limit.

#### Example 2

A slightly more complex example:

```yaml
domain: messaging
descriptors:
  # Only allow 5 marketing messages a day
  - key: message_type
    value: marketing
    descriptors:
      - key: to_number
        rate_limit:
          unit: day
          requests_per_unit: 5

  # Only allow 100 messages a day to any unique phone number
  - key: to_number
    rate_limit:
      unit: day
      requests_per_unit: 100
```

In the preceding example, the domain is "messaging" and we setup two different scenarios that illustrate more
complex functionality. First, we want to limit on marketing messages to a specific number. To enable this, we make
use of *nested descriptor lists.* The top level descriptor is ("message_type", "marketing"). However this descriptor
does not have a limit assigned so it's just a placeholder. Contained within this entry we have another descriptor list
that includes an entry with key "to_number". However, notice that no value is provided. This means that the service
will match against any value supplied for "to_number" and generate a unique limit. Thus, ("message_type", "marketing"),
("to_number", "2061111111") and ("message_type", "marketing"),("to_number", "2062222222") will each get 5 requests
per day.

The configuration also sets up another rule without a value. This one creates an overall limit for messages sent to
any particular number during a 1 day period. Thus, ("to_number", "2061111111") and ("to_number", "2062222222") both
get 100 requests per day.

When calling the rate limit service, the client can specify *multiple descriptors* to limit on in a single call. This
limits round trips and allows limiting on aggregate rule definitions. For example, using the preceding configuration,
the client could send this complete request (in pseudo IDL):

```
RateLimitRequest:
  domain: messaging
  descriptor: ("message_type", "marketing"),("to_number", "2061111111")
  descriptor: ("to_number", "2061111111")
```

And the service with rate limit against *all* matching rules and return an aggregate result.

#### Example 3

An example to illustrate matching order.

```yaml
domain: edge_proxy_per_ip
descriptors:
  - key: ip_address
    rate_limit:
      unit: second
      requests_per_unit: 10

  # Black list IP
  - key: ip_address
    value: 50.0.0.5
    rate_limit:
      unit: second
      requests_per_unit: 0
```

In the preceding example, we setup a generic rate limit for individual IP addresses. The architecture's edge proxy can
be configured to make a rate limit service call with the descriptor ("ip_address", "50.0.0.1") for example. This IP would
get 10 requests per second as
would any other IP. However, the configuration also contains a second configuration that explicitly defines a
value along with the same key. If the descriptor ("ip_address", "50.0.0.5") is received, the service will
*attempt the most specific match possible*. This means
the most specific descriptor at the same level as your request. Keep in mind that equally specific descriptors are matched on a first match basis. Thus, key/value is always attempted as a match before just key.

#### Example 4

The Ratelimit service matches requests to configuration entries with the same depth level. For instance, the following request:

```
RateLimitRequest:
  domain: example4
  descriptor: ("key", "value"),("subkey", "subvalue")
```

Would **not** match the following configuration even though the first descriptor in
the request matches the descriptor in the configuration.

```yaml
domain: example4
descriptors:
  - key: key
    value: value
    rate_limit:
      -  requests_per_unit: 300
         unit: second
```

However, it would match the following configuration:

```yaml
domain: example4
descriptors:
  - key: key
    value: value
    descriptors:
      - key: subkey      
        rate_limit:
          -  requests_per_unit: 300
             unit: second
```

### Descriptor list definition

Each configuration contains a top level descriptor list and potentially multiple nested lists beneath that. The format is:

```
domain: <unique domain ID>
descriptors:
  - key: <rule key: required>
    value: <rule value: optional>
	rate_limit: (optional block)
	  unit: <see below: required>
	  requests_per_unit: <see below: required>
    descriptors: (optional block)
	  ... (nested repitiion of above)
```

Each descriptor in a descriptor list must have a key. It can also optionally have a value to enable a more specific
match. The "rate_limit" block is optional and if present sets up an actual rate limit rule. See below for how the rule
is defined. The reason a rule might not be present is typically if a descriptor is a container for a 2nd level
descriptor list. Each descriptor can optionally contain a nested descriptor list that allows for more complex matches
and rate limit scenarios.

### Rate limit definition

```
rate_limit:
  unit: <second, minute, hour, day>
  requests_per_unit: <uint>
```

The rate limit block specifies the actual rate limit that will be used when there is a match.
Currently the service supports per second, minute, hour, and day limits. More types of limits may be added in the
future based on customer demand.

## Loading Configuration

The ratelimit service uses a library written by Lyft called goruntime to do configuration loading. Goruntime monitors
a designated path, and watches for symlink swaps to files in the directory tree to reload configuration files.

The path to watch can be configured via the [settings](https://github.com/lyft/ratelimit/blob/master/src/settings/settings.go)
package with the following environment variables:

```
RUNTIME_ROOT default:"/srv/runtime_data/current"`
RUNTIME_SUBDIRECTORY
```

For more information on how runtime works you can read its [README](https://github.com/lyft/goruntime).

# Rate limit statistics

The rate limit service generates various statistics for each configured rate limit rule that will be useful for end
users both for visibility and for setting alarms.

Rate Limit Statistic Path:

```
ratelimit.service.rate_limit.DOMAIN.KEY_VALUE.STAT
```

DOMAIN:
* As specified in the domain value in the YAML runtime file

KEY_VALUE:
* A combination of the key value
* Nested descriptors would be suffixed in the stats path

STAT:
* near_limit: Number of rule hits over the NearLimit ratio threshold (currently 80%) but under the threshold rate.
* over_limit: Number of rule hits exceeding the threshold rate
* total_hits: Number of rule hits in total

These are examples of generated stats for some configured rate limit rules from the above examples:

```
ratelimit.service.rate_limit.mongo_cps.database_default.over_limit: 0
ratelimit.service.rate_limit.mongo_cps.database_default.total_hits: 2846
ratelimit.service.rate_limit.mongo_cps.database_users.over_limit: 0
ratelimit.service.rate_limit.mongo_cps.database_users.total_hits: 2939
ratelimit.service.rate_limit.messaging.message_type_marketing.to_number.over_limit: 0
ratelimit.service.rate_limit.messaging.message_type_marketing.to_number.total_hits: 0
```

# Debug Port

The debug port can be used to interact with the running process.

```
$ curl 0:6070/
/debug/pprof/: root of various pprof endpoints. hit for help.
/rlconfig: print out the currently loaded configuration for debugging
/stats: print out stats
```

You can specify the debug port with the `DEBUG_PORT` environment variable. It defaults to `6070`.
