SDETool
=======

Uses the Dust514 **s**tatic **d**ata **e**xport to poll for info.

Highly unoptimized however it should be fairly stable.

Usage
=====

```
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
```
SDETool -i 364035
```
returns something like
```
Getting stats on Logistics ak.0
===== Description =====
The Logistics dropsuit is outfitted with the latest in integrated diagnostic technology, most of which revolves around maintaining the condition and efficiency of squad mates and their equipment. As such, a soldier equipped with this class of dropsuit becomes a force multiplier, greatly improving the overall effectiveness of the unit.

The Amarr variant is a durable, combat-focused suit that provides above-average protection, allowing logistic units to operate in the middle of a firefight, actively dispersing aid and support as needed while simultaneously engaging the enemy and inflicting trauma of its own.

When deployed, a soldier equipped with a Logistics suit fills a vital tactical role in small unit operations and full-scale warfare, providing both, medical and mechanical support.
-> Costs 57690 ISK
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

TODO
====
Use SQLite indexes for:
``` SQLite
"SELECT * FROM CatmaAttributes WHERE catmaAttributeName == 'mDisplayName' AND catmaValueText LIKE '%" + name + "%'"
```
Can take a really long time to get very vague things like "combat", "rail", etc..

More input types for -i use the typeName and display name to lookup with the best match working?
