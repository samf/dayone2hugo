# daypub: helping you make Day One journal exports public

This is a CLI tool that takes as input a zip file produced by [Day
One](https://dayoneapp.com). Its output is a set of files in a
directory of your choosing, which comprise the journal entry you
specified.

This is *almost* trivial, since Day One is good about not obfuscating
its output. But there are some minor nuisances. The most challenging aspect
is getting the images in the markdown to show up as links to the
images from the zip file. This tool fixes this automatically.

## installing

Two ways to install:

1. via [Homebrew](https://brew.sh/): `brew install samf/samf/daypub`
2. via go source:
    1. `git clone https://github.com/samf/daypub`
    2. (look at the source if you wish)
    3. `go build`
    4. `mv daypub <destination-in-your-$PATH>`

## using

Use `daypub --help` for help at any time.

### method 1: markdown

This makes an `index.md` file, with image links adjusted to the
correct file names of any images that are bundled.

### method 2: hugo

This is a superset of the `markdown` method. It adds the "front
matter" to the front of the `index.md` file, setting the title of
the entry if possible.

Optionally, it can use the `<figure>` shortcode instead of the
markdown image syntax. `<figure>` is especially useful if you intend
to keep your images in the page bundle along with your `index.md`.
