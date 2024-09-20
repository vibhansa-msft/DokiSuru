# DokiSuru

DokiSuru is a CLI utility to sync your local file to blockbloc. This is a delta sync utility so only the modified data on local disk will be uploaded.
As of now this utility works for a file and only syncs local data to remote. In future we will add support for directory and remote to local sync as well.


# Build
```
go build
```

# Usage
```
./dokisuru [OPTS]

  -path string
        Path to the file to process 

  -blocksize uint
        Block size to use (default 16777216)

  -validate
        Validate file after sync (default true). This will result in download of entire file from container.

  -worker int
        Number of parallel threads to use (default 16)
```

# Future Scope
- DokiSuru can sync only files which are ogirinally uploaded or synced earlier using the same utility.
- If DokiSuru figures that blob was modified outside the scope of this utility then it will fail to sync.
- If the file does not exists in container then entire file will be uploaded.
