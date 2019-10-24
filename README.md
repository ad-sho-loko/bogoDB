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

- go 1.13
- protoc

## How to run

```
# start bogodb server
> go run .

# create table
> curl "http://localhost:32198/execute?query=create%20table%20users%20{%20id%20int%20primary%20key%20%20}"

# insert 
> curl "http://localhost:32198/execute?query=insert%20into%20users%20values%20(1)"

# select
> curl "http://localhost:32198/execute?query=select%20id%20from%20users"
```

## TODO 

- refactoring `query`, especially analyse, eval...
- btree's implementation
- add update, delete statement

## Author

ad-sho-loko

## LICENSE

MIT