# dylib_pkg

Package dynamic libraries on MacOS.

Developed by [Blacknut](http://www.blacknut.com).

## Install

You must have [go](https://golang.org) installed, then:

```bash
   $ go install github.com/blacknut/dylib_pkg
```

## Usage

Just give your executable path:

```bash
  $ dylib_pkg /path/to/your/executable
```

All dynamics libs are then copied into to the same directory as your executable, and all references are modified to use `@executable_path`.

Get more help with:

```bash
  $ dylib_pkg -help
```

## References

- <http://stackoverflow.com/a/11585225>
