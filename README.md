# ccBatchAnalyzer utility (ccba)
This small console utility might be helpful for people who create scanword boards 
for Kryss (https://www.kryds.app/) using Crossword-Compiler (https://www.crossword-compiler.com/).

The utility reads all boards in the given directory (only XML files), looks for duplicate words,
and creates a simple CSV report which can be analysed later.

## Usage
"It's better to see something once..."
[![asciicast](https://asciinema.org/a/321964.svg)](https://asciinema.org/a/321964)

## Flags
* **--boards** Specifies path to the directory where all boards placed.
The program walks through directories recursively (which means that included directories will be processed)
* **--report** Specifies a path to the report file (by default a report is saved as _report.csv_ in the same directory)
* **--verbose** Prints info about each word to a console and other information which might be useful
* **--debug** Prints debug information to a console
