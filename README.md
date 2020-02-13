# repository-lister

Tool for listing public GitHub repositories.

## Considerations

* Max number of results: 1000 (because of GitHub API constraints)
* Requires a valid GitHub access token
* Most starred repositories first
* Included fields are: *repository name*, *stars*, *created at* and *description*
* Limited only to Golang repositories

## Usage

**repository-lister** requires a valid GitHub access token, with permissions to read public repositories.

```
./main -token=GITHUB_ACCESS_TOKEN
```

Results will be printed on the console.

## License

See the [LICENSE](LICENSE) file for license rights and limitations (MIT).

