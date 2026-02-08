# BackItUp

Simple CLI tool for backing up databases. Currently supports MongoDB, MySQL, and PostgreSQL.

## Installation

Using Make (recommended):
```bash
make build
```

Or using Go directly:
```bash
go build
```

Other available make targets:
- `make install` - Install dependencies
- `make test` - Run tests
- `make clean` - Remove build artifacts
- `make fmt` - Format code
- `make help` - Show all available targets

## Version

Check the installed version:
```bash
./BackItUp version
```

## Usage

### MongoDB

First time setup:

```bash
./BackItUp mongodb --config --uri "mongodb://localhost:27017"
```

Run backup:

```bash
./BackItUp mongodb
```

Backups go to `BACKUP/mongo/`

### MySQL

Configure your connection:

```bash
./BackItUp mysql --config --host localhost
./BackItUp mysql --config --user root
./BackItUp mysql --config --password yourpassword
./BackItUp mysql --config --database dbname
```

Run backup:

```bash
./BackItUp mysql
```

Backups go to `BACKUP/mysql/`

### PostgreSQL

Configure your connection:

```bash
./BackItUp postgresql --config --host localhost
./BackItUp postgresql --config --user postgres
./BackItUp postgresql --config --password yourpassword
./BackItUp postgresql --config --database dbname
```

Run backup:

```bash
./BackItUp postgresql
```

Backups go to `BACKUP/postgresql/`

## Config

Settings are saved in `config.yaml` in the current directory. You can also edit this file directly if you want.

---

https://roadmap.sh/projects/database-backup-utility
