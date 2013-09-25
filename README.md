# Integrity

## Compiling

    CGO_CFLAGS="-I`pg_config --includedir-server`" go build

## Running

    pg_integrity -f basebackup.tar
    pg_basebackup -D - -Ft | pg_integrity
