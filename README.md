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

## Configuration Status

Check your current configuration for all databases:
```bash
./BackItUp status
```

This displays:
- Connection settings for MongoDB, MySQL, and PostgreSQL
- Which databases are configured and ready to use
- Password masking for security
- Location of config file

## List Backups

View all available backups with sizes and timestamps:
```bash
./BackItUp list
```

This will show all backups organized by database type (MongoDB, MySQL, PostgreSQL) with their file sizes and creation dates.

## Restore from Backup

Restore a database from a previously created backup:

```bash
# Interactive mode - select from available backups
./BackItUp restore mongodb
./BackItUp restore mysql
./BackItUp restore postgresql

# Restore latest backup automatically
./BackItUp restore mysql --latest

# Restore specific backup file
./BackItUp restore postgresql --file BACKUP/postgresql/mydb.sql
```

The restore command will:
- Show available backups with timestamps
- Ask for confirmation before restoring
- Use the appropriate database tool (mongorestore, mysql, psql)

## Cleanup Old Backups

Manage backup storage with retention policies:

```bash
# Keep only last 30 days of MySQL backups
./BackItUp cleanup mysql --days 30

# Keep only 5 most recent PostgreSQL backups
./BackItUp cleanup postgresql --keep 5

# Clean up all databases at once
./BackItUp cleanup all --days 7

# Preview what would be deleted (dry run)
./BackItUp cleanup mysql --days 30 --dry-run
```

The cleanup command helps you:
- Save disk space by removing old backups
- Maintain retention policies
- Preview changes with `--dry-run` before actual deletion

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
