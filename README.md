# mdmdevsvc

Device information micro service for MicroMDM.

The aim of this project is to store information about devices under management by the MDM.
It is also responsible for periodically querying inventory information about those devices.

## Periodic Tasks

On startup:

- Perform a full fetch from DEP (batched)

At intervals (to be confirmed):

- Perform a DEP sync.
- Perform an Installed Application query.
- Perform a Device Information query.
- Perform an Installed Certificates query.
- Perform an OS Update Status check.
