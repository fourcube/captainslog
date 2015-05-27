# captainslog

Remember what you have done through a simple text log.

## Installation

Requires a valid go installation.

`go get -u github.com/fourcube/captainslog`

## Usage

Set $CAPTAINSLOG to the path of your logfile. Set $EDITOR if you don't want to use nano.

`captainslog` spawns your `$EDITOR` and shows you an output like:

```
|
## END ##
#########

# Date 2015-05-27 09:38:00 +0200 CEST
# ===================================
# I improved 'captainslog'.
# It should now be possible to add newlines and get a look at the
# last few entries.
# ...
# <Up to 5 previous entries>
```

You can then make your entry and simply close the file. It will be appended to the file at `$CAPTAINSLOG`.


