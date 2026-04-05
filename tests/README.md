1.  Bind to a port:
2.  Respond to PING:

```shell
nc localhost 6379
Connecting to port 6379...
redis-cli PING
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
redis-cli ECHO banana
```

6.  Implement the SET & GET commands

```shell
redis-cli SET foo bar
redis-cli GET foo

redis-cli SET mango orange
GET mango
```

7. Expiry:

```shell
redis-cli SET foo bar PX 2000
redis-cli GET foo

redis-cli SET mango pineapple PX 100
GET mango
GET mango
```

8.  The INCR command:

```shell
redis-cli SET blueberry 41
INCR blueberry

redis-cli INCR banana
INCR banana
GET banana

redis-cli SET raspberry grape
INCR raspberry
```

9. The MULTI and EXEC command:

```shell
redis-cli MULTI

redis-cli EXEC

redis-cli MULTI
EXEC
EXEC

redis-cli MULTI
SET pear 98
INCR pear
redis-cli GET pear

redis-cli MULTI
SET blueberry 40
INCR blueberry
INCR grape
GET grape
EXEC
redis-cli GET blueberry
```
