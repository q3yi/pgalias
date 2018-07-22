# Alias for `psql` and `pg_dump` management

Save all your Postgres connection with its alias and connect database using
`psql` or `pg_dump`

## Config file

A bit like `.pgpass` file but separated using <kbd>TAB</kbd>, and the first
element is the alias for current connection. The default alias filepath is
`~/.pgapass`

```text
local	localhost	5432	postgres_db	postgres_user	postgres_password
```
