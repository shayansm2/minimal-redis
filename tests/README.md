1.  Bind to a port:
2.  Respond to PING:

```shell
nc localhost 6379
```

3.  Respond to multiple PINGs:
4.  Handle concurrent clients:

```shell
redis-cli PING
echo -e "PING\nPING" | redis-cli
nc localhost 6379
```
