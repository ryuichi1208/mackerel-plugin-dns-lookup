# mackerel-plugin-dns-lookup

This tool is a Mackerel plug-in that displays the time it takes to send a specific query to a specified DNS server for a specified number of times.There are three types of metrics that are output: "Max", "Avg", and "Min".

## Usage

```
Usage:
  main [OPTIONS]

Application Options:
  -s, --server=   DNS Server (default: 8.8.8.8)
  -p, --port=     query num (default: 53)
  -d, --domain=   Domain Name
  -t, --type=     Record Type
  -n, --count=    query num (default: 3)
      --timeout=  deadline time
      --protocol= Record Type (default: udp)
      --threads=  thread num (default: 1)
      --version   show version
      --debug

Help Options:
  -h, --help      Show this help message

Usage:
  main [OPTIONS]

Application Options:
  -s, --server=   DNS Server (default: 8.8.8.8)
  -p, --port=     query num (default: 53)
  -d, --domain=   Domain Name
  -t, --type=     Record Type
  -n, --count=    query num (default: 3)
      --timeout=  deadline time
      --protocol= Record Type (default: udp)
      --threads=  thread num (default: 1)
      --version   show version
      --debug

Help Options:
  -h, --help      Show this help message
```

Run (Unit is ms)

```
$ mackerel-plugin-dns-lookup \
    --server 8.8.8.8 \
    --port 53 \
    --domain google.com \
    --count 10 \
    --threads 1

min	10	1663065111
max	23	1663065111
avg	12	1663065111
```
