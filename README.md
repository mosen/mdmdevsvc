# mdmdevsvc

Device information service for [microMDM](https://micromdm.io). *Located on GitHub [here](https://github.com/micromdm/micromdm)*

This application stores information about all the devices managed by the MDM. Features:

- Discover new devices via DEP.
- Periodically query and store information about device attributes, applications, certificates, and settings.
- Provide an API to fetch existing device information or to refresh the information if another service knows that
it is out of date.

## Configuration

The database is assumed to be PostgreSQL.

You have several options for providing configuration: by command line flag, environment variable, or toml.


The example configuration file is [here](mdmdevsvc.toml.example)

*NOTE:* This service uses [flaeg](https://github.com/containous/flaeg) and [staert](https://github.com/containous/staert) packages
for runtime configuration.

## Periodic Tasks

The device service needs to keep information about changes in devices up to date, which means there's a scheduled task
for querying your devices for varying kinds of information.

These tasks are performed automatically, and you can configure how often they happen:

On startup:

- Perform a full fetch from DEP (batched)

At intervals:

- Perform a DEP sync (check for new or changed devices since the last DEP fetch).
- Perform an Installed Application query.
- Perform a Device Information query.
- Perform an Installed Certificates query.
- Perform an OS Update Status check.

## Ideas

- Webhooks/Callbacks on cert/app install finished and confirmed.
