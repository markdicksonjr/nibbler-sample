# Local Auth Sample

## Building

You can build and run the old-fashioned way from this directory:

```bash
go run main.go
```

Or, you can use nibbler's build utility with

```bash
go run build/main.go
```

Which will build multiple platforms in the ./dist directory, and embed the git tag into the binaries.
This sample app prints the git tag when it is ran.