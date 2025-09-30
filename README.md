# daypub: helping you make Day One journal exports public

This is a CLI tool that takes as input a zip file produced by [Day
One](https://dayoneapp.com). Its output is a set of files in a
directory of your choosing, which comprise the journal entry you
specified.

This is *almost* trivial, since Day One is good about not obfuscating
its output. But there are some minor nuisances. The hardest problem
is getting the images in the markdown to show up as links to the
images from the zip file.

## installing

Two ways to install:

1. via [homebrew](https://brew.sh/): `brew install samf/samf/daypub`
2. via go source:
  1. `git clone https://github.com/samf/daypub`
  2. (look at the source if you wish)
  3. `go build`
  4. `mv daypub <destination-in-your-$PATH>`
