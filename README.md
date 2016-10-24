# rello

Track actions over time with trello check lists.

### Usage

Point a trello webhook for a checklist to wherever this is running.
[Trello Api Docs](https://developers.trello.com/)

Create the sqlite3 database with the init script.
```
$ ./init.sh
```

Export port you want this to run on and the path to sqlite3 database.
```
export PORT=5040
export SQLITE3=$GOPATH/src/github.com/adamryman/rello/db/db.db

```

Install and Run.
```
$ go install github.com/adamryman/rello/cmd/rello
$ rello

```
