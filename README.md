# Uptime Monitor

This is a simple application that tracks website uptime and reports total uptime and outages.

## Requirements

Originally this was intended to be a stand alone app that would run on its own. SQLite would be used to make this super portable and easy, but that was turning into a headache, so it was swapped out for MariaDB/MySQL. To get your database in place you can run the sql migration files in the migrations folder.

The migrations are written in plain SQL, but you can use the awesome tool [`goose`](https://github.com/pressly/goose/) to run the migrations and keep your DB in sync.
