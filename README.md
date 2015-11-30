# THIS IS DEPRECATED
Please use https://github.com/prometheus/blackbox_exporter instead.

# Blackbox Prober

[![Build Status](https://travis-ci.org/jswank/blackbox_prober.svg)](https://travis-ci.org/jswank/blackbox_prober)

Export blackbox telemetry like availability, request latencies and
request size for remote services.

## Supported URLs
### http/https
The exporter requests the given url and reads from it until EOF.

### tcp
The exporter connects to the given host:port. If any path is given, it
will try to read until EOF which is required for exposing the size.

### icmp
Execute `ping`. Port and path are ignored.

## Available metrics

### All
- `blackbox_up{url}` 1 if url is reachable, 0 if not
- `blackbox_latency_seconds{url}` Latency of request for url

### http/https
- `blackbox_size_bytes{url}` Total size of response for url 
- `blackbox_cert_expire_timestamp{url}` Expiry date of certificate (HTTPS only)
- `blackbox_response_code{url}` Status code for the URL

## tcp
- `blackbox_size_bytes{url}` Size of the response

## Example

    ./blackbox_prober \
      -u http://5pi.de \
      -u https://5pi.de \
      -u icmp://192.168.178.1 \
      -u tcp://freigeist.org:655

## Using Docker

    docker pull jswank/blackbox_prober

    docker run -d -p 9110:9110 jswank/blackbox_prober \
        -u http://5pi.de \
        -u https://5pi.de \
        -u icmp://192.168.178.1 \
        -u tcp://freigeist.org:655
