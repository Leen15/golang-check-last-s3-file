# Check last S3 file
A useful script that search for last file in a s3 bucket and check size and last edit timestamp.
We use this with a cron job for check that our backups are up to date.

Following environment variables are mandatory:
```
TIMEZONE=                  // ex. Europe/Rome
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=                // ex. eu-west-1
BUCKET=                    // Only the bucket name ex. Backups
PATH_PREFIX=               // You can insert here subpath and the name prefix ex. /databases/my-database
OLDER_THAN_MINUTES=        // ex. 60 for 1 hour
SMALLER_THAN_MEGABYTES=    // ex. 100 for 100MB
```

## Usage

`docker run --env-file .env leen15/check-last-s3-file`

And it responses with something like:

```
Check last file in:
Bucket:          backups
Prefix:          databases/my-database
Older than:      60 minutes
Smaller than:    100 MB
Files in bucket: 4

Most recent file is:
Name:          my-database.2019-11-22-00-10-09.native.dump.bz2
modified at:   2019-11-22 01:10:11 +0100 CET
Size:          50 MB

[ERROR] The file is older than max allowed. ( 13h57m23.59s ago )
[ERROR] The file is smaller than min allowed. ( 50 vs 100 MB )
```

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
