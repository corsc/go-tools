This tool intents to calculate the test coverage of a particular package (including any sub packages).

## Commands:

* `$ package-coverage -c ./` will generate coverage.  1 coverage file (*.cov) per package
* `$ package-coverage -d ./` will remove any previous coverage files (will remove all *.cov files)
* `$ package-coverage -p ./` will import all coverage files under the supplied dir and output the summary.
* `$ package-coverage -c -p -d ./` all of the above

## Output Sample
```
  %		Statements	Package
50.00		  240	sage42.org/go-tools/package-coverage/
42.47		   73	sage42.org/go-tools/package-coverage/generator/
65.52		   87	sage42.org/go-tools/package-coverage/parser/
58.18		   55	sage42.org/go-tools/package-coverage/utils/
```
