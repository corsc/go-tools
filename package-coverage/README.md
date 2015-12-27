This tool intents to calculate the test coverage of a particular package (including any sub packages).

To use this tool first run the "coverage.sh" script in the base go package that you want to calculate coverage for.

Once this is complete there should be an acc.out file created in the directory the command was run.

Then you can run:
    $ ./package-coverage acc.out

You should receive an output similar to:

       %		Statements	Package
     63.22		   87		sage42.org/go-tools/
     63.22		   87		sage42.org/go-tools/package-coverage/
