# go-dupfinder

We are currently working on a tool that will find duplicate files in a directory.

Note: Always use latest version

## WIP : use it at your own risk

Work in progress.

I will complete this project in my spare time later.

## Building

For Production,

```sh
go build -ldflags="-X 'github.com/mrinjamul/go-dupfinder/app.Version=$(git describe --tags $(git rev-list --tags --max-count=1) || echo "dev")' -X 'github.com/mrinjamul/go-dupfinder/app.GitCommit=$(git rev-parse HEAD)'"
```

## License

open-sourced under [MIT License](LICENSE)
