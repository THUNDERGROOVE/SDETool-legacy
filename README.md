SDETool
=======

Uses the Dust514 __S__tatic __D__ata __E__xport to poll for info.

Highly unoptimized for now however it should be fairly stable(hopefully).


Building
========
You need [Go](http://golang.org) with your GOPATH environment variable setup our dependencies

On Windows you will need a matching version of GCC which you can usually get from either [Cygwin](http://www.cygwin.com/) or [Mingw](http://www.mingw.org/)  I use Mingw so that's your best bet.  It must match the CPU architecture as what you're building for (386, amd64).

In Linux you'll need build-essentials.  Installing it can depend on your distro however you should be able to figure this out if you need.
``` bash
make dep
```
Should download our dependencies for you

Then you should be able to
``` bash
make
```
And have your SDETool.exe or SDETool binary!

If you have any issues building just submit an issue and I'll do my best sorting it out however it should work on any officially supported platform from Go(Windows, Linux, FreeBSD, OpenBSD and Plan9).  Currently I test it in Windows 8.1 and Debian testing.

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

Along with exact values you can also get tables of values by:
``` bash
SDETool -d "Kaalakiota rail rifle"
```
Outputs:
``` 
Getting damage on: Kaalakiota Rail Rifle
Damage mods[cmplx]|Proficiency level |Output damage     
0                 |0                 |60                
0                 |1                 |62                
0                 |2                 |63                
0                 |3                 |65                
0                 |4                 |66                
0                 |5                 |68                
1                 |0                 |65                
1                 |1                 |67                
1                 |2                 |69                
1                 |3                 |70                
1                 |4                 |72                
1                 |5                 |74                
2                 |0                 |71                
2                 |1                 |73                
2                 |2                 |75                
2                 |3                 |77                
2                 |4                 |79                
2                 |5                 |81                
3                 |0                 |77                
3                 |1                 |79                
3                 |2                 |81                
3                 |3                 |83                
3                 |4                 |85                
3                 |5                 |87                
4                 |0                 |80                
4                 |1                 |82                
4                 |2                 |84                
4                 |3                 |86                
4                 |4                 |88                
4                 |5                 |91                
5                 |0                 |82                
5                 |1                 |84                
5                 |2                 |86                
5                 |3                 |88                
5                 |4                 |90                
5                 |5                 |92               
```

What works?
===========
1. Damage calculations (With generic damage tables)
2. Searching
3. Dynamic info printing
4. Simple gathering of skill bonus (Currently only works on additive skills that don't work off of a percentage)
5. Recusive skill lookups. I.E. it knows it takes Amarr medium dropsuits level 3 to use Amarr assault dropsuits.
6. 

TODO
====

1. <del>More calculations like damage for things like dampening, range amps, speed, etc.</del>  (Implemented GenericCalculateValue, should make it easier when I get to it)
2. Fix math for skill applications, ocasionally it'll break
3. Use an application directory for DB, log, zips and.  %APPDATA% for windows and .SDETool for linux/OSX