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

5.  Implement the ECHO command

```shell
redis-cli PING
redis-cli ECHO hey
echo -e "PING\nECHO hi" | redis-cli
```

6.  Implement the SET & GET commands

```shell
redis-cli SET foo bar
redis-cli GET foo
```
