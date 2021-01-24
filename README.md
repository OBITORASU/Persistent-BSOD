# Persistent-BSOD
A BSOD program that registers itself as a startup and windows service.

# How it works
The code springs into action as soon as a user is tricked into running the .scr file on their Windows machine. Upon execution, the program tries to access restricted memory which causes Windows to trigger an access violation resulting in a BSOD. To maintain persistence on the infected PC, the program makes changes to the Windows Registry so that it is run at startup on each subsequent reboot.

# Try it out
Get it here on the [releases](https://github.com/No-Cellist-7780/Persistent-BSOD/releases/tag/1.0.0) page.

