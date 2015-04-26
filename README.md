# Ping Exporter
Export availability, request latencies and size for remote services.

## Supported URLs
### http/https
The exporter requests the given url and reads from it until EOF.

### tcp
The exporter connects to the given host:port. If any path is given, it
will try to read until EOF which is required for exposing the size.

### icmp
Execute `ping`. Port and path are ignored.

## Available metrics
- `ping_up{url}` 1 if url is reachable, 0 if not
- `ping_latency_seconds{url}` Latency of request for url
- `ping_size_bytes{url}` Size of request for url

## Example

    ./ping_exporter \
      -u http://5pi.de \
      -u https://5pi.de \
      -u icmp://192.168.178.1 \
      -u tcp://freigeist.org:655

