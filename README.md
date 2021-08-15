# txtimg

Util to covert text to PNG image.

# install

```
go install github.com/epii1/txtimg
```

# usage

```
cat file.txt | txtimg -fontfile path/to/font/file > out.png
```

```
Usage of ./txtimg:
  -chars int
        chars displayed per line (default 20)
  -dpi float
        screen resolution in Dots Per Inch (default 240)
  -fontfile string
        filename of the ttf font
  -height int
        image height in points (default 400)
  -hinting string
        none | full (default "none")
  -padding int
        text left and right padding (default 10)
  -size float
        font size in points (default 14)
  -spacing float
        line spacing (default 1)
  -whiteonblack
        white text on a black background
  -width int
        image width in points (default 940)
```
