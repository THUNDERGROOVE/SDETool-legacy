SDETool
=======

Uses the Dust514 **S**tatic **D**ata **E**xport to poll for info.

Highly unoptimized for now however it should be fairly stable(hopefully).

[![Gobuild Download](http://gobuild.io/badge/github.com/THUNDERGROOVE/SDETool/download.png)](http://gobuild.io/github.com/THUNDERGROOVE/SDETool)

Building
========
You need [Go](http://golang.org) with your GOPATH environment variable setup our dependencies
```
go get github.com/mattn/go-sqlite3
```
Then you should be able to
```
make
```
And have your SDETool.exe or SDETool binary!

If you have any issues building just submit an issue and I'll do my best sorting it out however it should work on any officially supported platform.  Currently I test it in Windows 8.1 and Debian testing.

Usage
=====

``` bash
SDETool -s "Logistics ak.0"
```

returns

```
Searching value: 'Logistics ak.0'
365308 | 'Neo' Logistics ak.0
365714 | 'Pyrus' Logistics ak.0
366420 | Imperial Logistics ak.0
364035 | Logistics ak.0
```

useful for getting typeIDs

``` bash
SDETool -i 364035
# or
SDETool -i "Logistics ak.0"
```

returns something like

```
Getting stats on Logistics ak.0
===== Description =====
The Logistics dropsuit is outfitted with the latest in integrated diagnostic technology, most of which revolves around maintaining the condition and efficiency of squad mates and their equipment. As such, a soldier equipped with this class of dropsuit becomes a force multiplier, greatly improving the overall effectiveness of the unit.

The Amarr variant is a durable, combat-focused suit that provides above-average protection, allowing logistic units to operate in the middle of a firefight, actively dispersing aid and support as needed while simultaneously engaging the enemy and inflicting trauma of its own.

When deployed, a soldier equipped with a Logistics suit fills a vital tactical role in small unit operations and full-scale warfare, providing both, medical and mechanical support.
-> Cost 57690 ISK
===== Dropsuit =====
-> Heavy Weapons: 0
-> Light Weapons: 1
-> Sidearms: 1
-> Equipment slots: 3
-> High slots: 3
-> Low slots: 4
===== Tags =====
-> 352332 tag_dropsuit
-> 352339 tag_amarr
-> 353508 tag_dropsuit_logistics
-> 353502 tag_core
```

Currently, we only support getting damage values for weapons that don't get their damage inherited from their projectile type.  An example would be to:
``` bash
# Get damage for a Kaalakiota Rail Rifle with 3 complex damage mods and proficiency level 3
SDETool -d 365448 -c 2 -p 3
# You can also do things like
SDETool -d "kaalakiota rail rifle" -c 2 -p 3
# or
SDETool -d wpn_railrifle_ca_pro -c 2 -p 3
```

It will display a 64 bit float which should go out to 16(?) decimal places.  _This may change_

TODO
====

1. More calculations like damage for things like dampening, range amps, speed, etc.
..1. (Implemented GenericCalculateValue, should make it easier when I get to it)
2. <del>I'm currently working on implementing a faster way of using the -s flag.  
.. It takes about 0.1 seconds for each type that matches our search pattern so when we have > 30 items it takes a long time to return a result.  Will be using fewer SQL querries and possible GoRoutines and channels to get results back faster.</del> 
.. Mostly done, would still like faster searches