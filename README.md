# Hallo
Hallo is a side project to see if I can make a Bonjour-like program for an old laptop.

## Philosophy
An alarm clock can never really be a personal assistant. To make sure it stays a fancy alarm clock, Hallo operates on the philosophy that once the user starts the program, they should very rarely be forced to directly interact with it. Therefore, unlike Bonjour, Hallo doesn't feature any manner of voice recognition. 

## Preparing for install
You will have to install the portaudio development headers for any audio output to work. This should be quite easy on Linux. Precompiled headers are available for Windows from various sources.

If you want to build from source, you have to `go get` gocui (`go get github.com/jroimartin/gocui`). Glide DOES NOT WORK for gocui.
