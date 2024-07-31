# Test Paging of MySQL

Two ways of paging:

- `pageNo` and `pageSize`
- `lastId` and `pageSize`

## Environment

- Macbook Air M1
- MySQL 5.7
- Golang 1.22.4
- GORM

Test paging querying 1,000,000 users.

## Result

`go run . ${mysql_user} ${mysql_pass}`

```
-------------Paging by offset-------------
Paging by offset took:
135.720125ms
-------------Paging by lastId-------------
Paging by lastId took:
1.144125ms
```

