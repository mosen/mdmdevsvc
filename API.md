# DeviceStore REST API #

## Overview ##

The devicestore REST API conforms to the JSON-API v1 specification.

## Supported Endpoints ##

### /v1/devices ###

    GET /v1/devices

Get an index of devices

    POST /v1/devices
    
Create a new device record.

    PUT /v1/devices/{uuid}
    
Update an existing device record.

    DELETE /v1/devices/{uuid}
    
Delete an existing device record.

### /v1/groups ###

    GET /v1/groups
    
Get a list of groups

    POST /v1/groups
    
Create a new group record.

    PUT /v1/groups/{uuid}

Update a group record

    DELETE /v1/groups/{uuid}
    
Delete an existing group record

    GET /v1/groups/{uuid}/devices
    
Get devices that are part of the specified group.

### /v1/subnet/(network)/(mask)

    GET /v1/subnet/192.168.0.0/24

Get a list of devices that have interfaces in the specified net/mask range

### /v1/locations

    GET /v1/locations

Get a list of physical locations which may be associated with devices.

    POST /v1/locations
    
Create a new physical location record.

    PUT /v1/locations/{uuid}
    
Update an existing physical location record.

    DELETE /v1/locations/{uuid}
    
Delete a physical location record. Devices in that location will not have any physical
location associated with them after the delete.

    GET /v1/locations/{uuid}/devices
    
Get a list of devices in the specified physical location.

### /v1/by_attribute

Retrieve device lists by attribute.

    GET /v1/by_attribute/ProductName/iPad4,1
    

