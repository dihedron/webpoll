# webpoll

A very simple HTTP client to check if/when a URL is available by polling.

## Usage

To run this utility:

```bash
$ > webpoll --attempts 100 --timeout 200 http://www.example.com
```

It will attempt to connect up to 100 times to the given URL, timing out after 200 ms if the connection takes too long and ensuring that at least tantamount milliseconds elapse between successive attempts; as soon as the connection succeeds (HTTP status code in the 2xx or 3xx classes) it exits with ```$?``` equal to ```0```. If all attempts fail, it exits with ```1```.

