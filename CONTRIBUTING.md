# Contributing

Hey there! We're thrilled that you'd like to contribute to this project. Your help is essential for keeping it great.

Contributions to this project are released to the public under the project's open-source license.

Please note that this project is released with a [Contributor Code of Conduct](https://github.com/ziplinesci/ziplines-ci-foundation/blob/main/CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.


## How to Contribute

```shell
git clone github.com/ziplinesci/ziplines-ci-foundation
go get ./...
```

### Submitting Pull Requests
1. Fork(https://github.com/ziplinesci/ziplines-ci-foundation/fork) the repository and create your branch from main.
2. Configure and install the dependencies. `go get ./...`
3. Create a new branch: git checkout -b my-branch-name.
4. Make your change, add tests (if the logic are complex and have side effect), and make sure the format and lint tests pass. `go fmt ./... && go vet ./... && golangci-lint run --fix`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to your fork and submit a pull request
7. Pat yourself on the back and wait for your pull request to be reviewed and merged.
8. The maintainers will review your pull request and either merge it, request changes to it, or close it with an explanation. 

##  Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests)
