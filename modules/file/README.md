# File module

## Description

Module to manipulate files.  This provides the base functionality used by several other modules, including Template.

## Configuration

Parameter | Required | Default | Choices | Description
----------|----------|---------|---------|------------
path      | yes      |         |         | Path to the file
action    | no       | touch   | touch   | Touch a file
          |          |         | set     | Set the content of the file
          |          |         | append  | Append to the content of a file
          |          |         | remove  | Delete the file
content   | no       |         |         | Data to write to the file in set/append modes
mode      | no       | 0644    |         | File mode
context   | no       |         |         | Set a specific selinux context
owner     | no       |         |         | Set file owner
group     | no       |         |         | Set file group


## Development notes

In addition to string and []byte, content supports the pipe type, to allow large files to be passed efficiently.