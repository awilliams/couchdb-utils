couchdb-utils
=============

A fast, lightweight, and portable CouchDB utility. See help below for more information. Built with Go.

### Download

[Binaries for select systems are available](https://github.com/awilliams/couchdb-utils/releases)

There are no dependencies besides the included binary.

**Example Usage**
```bash
# refresh views in `mydb` database on host couch.example.com:1234
couchdb-utils refreshviews mydb --host=user:pass@couch.example.com:1234 -v

# start continuous replication of all databases on remote host `33.33.33.10:5984`
couchdb-utils rep host 33.33.33.10:5984 -h user:secret@33.33.33.11:5984 -v
```

**Base commands**
```bash
Usage:
  couchdb-utils [flags]
  couchdb-utils [command]

Available Commands:
  version                            :: Prints the version number of couchdb-utils
  server                             :: Print basic server info
  stats [(<part1> <part2>)]          :: Print server stats (optionally only a certain section eg: couchdb request_time).
  activetasks [<type>]               :: Print active tasks (optionally filtering by type)
  session                            :: Print information about authenticated user
  databases                          :: Print all databases
  views [<db>...]                    :: Print all views (optionally filtering by database(s))
  refreshviews [<db>...] [--verbose] :: Refresh views (optionally filtering by database(s))
  rep <command>...                   :: Replication subcommands
  help [command]                     :: Help about any command

 Available Flags:
  -d, --debug=false: print http requests
  -h, --host="http://localhost:5984": Couchdb server url (http://user:password@host:port)
  -v, --verbose=false: chatty output
```

**Replication commands**
```bash
Usage:
  couchdb-utils rep <command>... [flags]
  couchdb-utils rep [command]

Available Commands:
  list                                                 :: Print all replicators
  start <source> <target> [--create --continuous]      :: Configure replication from source to target
  stop (<id>... | --all) [--verbose]                   :: Stop replicating given id(s) or all
  host <remote_host> [--create --continuous --verbose] :: Replicates all databases in remote host

 Available Flags:
  -d, --debug=false: print http requests
  -h, --host="http://localhost:5984": Couchdb server url (http://user:password@host:port)
  -v, --verbose=false: chatty output
```

##### Compiling

A simple makefile is provided. Make sure GO is installed and setup for cross-compiling. See [here](http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go) and [here](https://coderwall.com/p/pnfwxg) for help.
