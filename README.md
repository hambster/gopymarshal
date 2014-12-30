go-pymarshal
============

Marshal and unmarshall data between golang and python

# Versions

## 0.1
* unmarshal python marshaled data
* marshal golang data as python marshaled data

# Notes
* for unmarshal, supported types are listed below:
    * None
    * int
    * string
    * unicode
    * float
    * list
    * dict
* for marshal, supported types are listed below:
    * nil
    * int32
    * float
    * string
    * slice (with the types above)
    * map (with the types above)
