# BackItUp

Simple CLI tool for backing up databases. Currently supports MongoDB and MySQL.

## Installation

```bash
go build
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

## Config

Settings are saved in `config.yaml` in the current directory. You can also edit this file directly if you want.
