# BogoDb

BogoDb is a toy database management system written in Go.
Inspired by CMU Database Group's Lecture (of course not homework!) This db still works completely poor. I realized it is so difficult to develop efficient database...

## Specification 

- SQL
    - create table
    - insert statment
    - select statement(from, where)
    - begin, commit, rollback
- Index(with b-tree)
- Buffer on memory
- Concurrency(only transaction)
- not mmap implementation

## Requirement

- go
- maybe unix 
- protoc

## How to run

```
# start bogodb server
> go run .

# create table
> curl `create table users { id int primary key }

# insert 
> curl `insert into users values (1)`

# select
> curl `select id from users`
```

## TODO 

- refactoring `query`, especially analyse, eval...
- btree's implementation
- add update, delete statement

## Author

ad-sho-loko

## LICENSE

MIT