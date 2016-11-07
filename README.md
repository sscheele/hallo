# Hallo
Hallo is a side project to see if I can make Bonjour for an old 
laptop.

## Preparing for install
Installation requires you to install the stuff you need for CMU 
Sphinx:

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

The paths to pocketsphinx and sphinxbase are the paths which 
contain `pocketsphinx.pc` and `sphinxbase.pc`.
