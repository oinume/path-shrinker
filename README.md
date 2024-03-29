[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE)
[![CircleCI](https://img.shields.io/circleci/build/github/oinume/path-shrinker/master.svg?style=for-the-badge)](https://circleci.com/gh/oinume/path-shrinker/tree/master)
[![Codecov branch](https://img.shields.io/codecov/c/github/oinume/path-shrinker/master.svg?style=for-the-badge)](https://codecov.io/gh/oinume/path-shrinker)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/oinume/path-shrinker)

# path-shrinker

path-shrinker is a command to shrink directory path like this.

```shell
$ path-shrinker -short -last -tilde /Users/oinuma/go/src/github.com/oinume/path-shrinker
~/g/s/g/o/path-shrinker
```

## Install

### Download a binary

You can download a binary from [GitHub](https://github.com/oinume/path-shrinker/releases).

## Customize your Bash prompt with path-shrinker

Define `PS1` in your .bashrc or .bash_profile.

```shell
PS1='$(path-shrinker -fish) $ '
```

Then your terminal will show a prompt like this:
```shell
~/g/s/g/o/path-shrinker $ pwd
/Users/kazuhiro/go/src/github.com/oinume/path-shrinker
```

## Examples

For this directory tree:

```
/home/
  me/
    foo/
      bar/
        quux/
      biz/     # The prefix b is ambiguous between bar and biz.
```

Here are the results of calling `path-shrinker <option> /home/me/foo/bar/quux`:

```
Option        Result
<none>        /h/m/f/ba/q
-l|--last     /h/m/f/ba/quux
-s|--short    /h/m/f/b/q
-t|--tilde    ~/f/ba/q
-f|--fish     ~/f/b/quux
```
