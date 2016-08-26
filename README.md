# SChange

SChange is basically a safer way to do chown chgrp and chmod for sudoers.



## The Problem

The problem is with the way the unix tools work by default. For example: `%somegroup ALL=(root) chown <user> /allowed/path/` Will only allow sudo permissions for the specific path. Although allowing special regex might be possible, it's still dangerous. So, you're a master at regex, and no one could do anything like `sudo chown someuser /allowed/path/legitfile /etc/shadow` or `sudo chown someuser /allowed/path/legitfile{, /etc/shadow}`. Cool. You're awesome, because regex in sudoers sucks. But what about the symlink `/allowd/path/legitfile -> /etc/shadow`. Chown will happily change the permissions on `/etc/shadow` for you. You can use the -h option at that point, sure. But at the end of the day, this was just fun to write, and solves a problem without worrying about making bullet-proof sudoers regex, without worrying about symlinks. And, it's fun. 



## The Fix

The fix is as complicated as the  problem. We need to test for bad things, and only allow the good. The safer change way is to use a root readable config file that holds some of our data in it, like allowed paths.



## They Why's

* Why a config file?
* Why re-write these awesome unix tools?
* Why go?
