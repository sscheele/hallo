# Hallo
Hallo is a side project to see if I can make a Bonjour-like program for an old laptop.

## Philosophy
An alarm clock can never really be a personal assistant. To make sure it stays a fancy alarm clock, Hallo operates on the philosophy that once the user starts the program, they should very rarely be forced to directly interact with it. Therefore, unlike Bonjour, Hallo doesn't feature any manner of voice recognition. 

## Preparing for install
Installation requires you to install the stuff you need for CMU Sphinx:

1. [Sphinxbase](https://github.com/cmusphinx/sphinxbase)
2. [Pocketsphinx](https://github.com/cmusphinx/pocketsphinx)

Both installs can be run like so:

```
./autogen.sh
make
sudo make install
```

Finally, grab Sphinx:

```
$ export PKG_CONFIG_PATH=path/to/sphinxbase:path/to/pocketsphinx
$ go get github.com/xlab/pocketsphinx-go/sphinx
```

Before running the program, you'll need to also include your 
libraries in the path:

```export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/localib```

The paths to pocketsphinx and sphinxbase are the paths which 
contain `pocketsphinx.pc` and `sphinxbase.pc`.
