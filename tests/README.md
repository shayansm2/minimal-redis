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

7. Expiry:

```shell
redis-cli SET foo bar PX 2000
redis-cli GET foo
```

8.  The INCR command:

```shell
redis-cli SET foo 5
redis-cli INCR foo
redis-cli INCR foo

redis-cli INCR missing_key
redis-cli GET missing_key

redis-cli SET foo xyz
redis-cli INCR foo
```

9. The MULTI command:

```shell
redis-cli MULTI
redis-cli EXEC
```
