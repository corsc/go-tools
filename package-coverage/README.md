This tool intents to calculate the test coverage of a particular package (including any sub packages).

## Command line options:

* `$ package-coverage -c ./` will generate coverage.  1 coverage file (*.cov) per package
* `$ package-coverage -d ./` will remove any previous coverage files (will remove all *.cov files)
* `$ package-coverage -p ./` will import all coverage files under the supplied dir and output the summary.
* `$ package-coverage -v` is useful for debugging as it will print to std out a trace of what it is doing
* `$ package-coverage -s` will switch this tool into "single directory" mode (will not recurse down the file tree)
* `$ package-coverage -slack -webhook=https://hooks.slack.com/services/fu/bar` will print the coverage information to Slack using the supplied webbook
* `$ package-coverage -i="_generated/"` defines a regex of directories that should be excluded from coverage (useful for generated code)
* `$ package-coverage -p -prefix="github.com/corsc/"` this string will removed from the front of any outputted package names (current only supported by the slack output)
* `$ package-coverage -slack -depth=1` how many levels to output.  This does not effect the calculation only the output. (current only supported by the slack output)
* `$ package-coverage -p -m=1` will highlight (in red) the console output of any packages below the supplied number (current only supported console output)

## Notes:
* The coverage and statements are recursive (except in single dir mode).  Meaning the values for ./packageA/ include the values from ./packageA/packageB/
* In order to calculate coverage for directories with no tests, this tool will make a fake test file called `fake_test.go` prior to running coverage calcuation.  It will also remove it when the calculation is complete.  Cancelling this tool mid-run could cause this file to be remain.  This file can be deleted.
* This tool is not smart enough to detect existing `fake_test.go` files and not remove them.  All such files will be deleted.
* If things don't look right, please run in verbose mode `-v` and include that in any bug report.

## Output Sample
```
  %		Statements	Package
50.00		  240	github.com/corsc/go-tools/package-coverage/
42.47		   73	github.com/corsc/go-tools/package-coverage/generator/
65.52		   87	github.com/corsc/go-tools/package-coverage/parser/
58.18		   55	github.com/corsc/go-tools/package-coverage/utils/
```
