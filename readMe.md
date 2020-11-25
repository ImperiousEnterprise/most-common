# Most Common

Most Common finds the three most common words in text and displays an ordered list by most frequent.
Data can come through piping or as a command line variable.

1. The program accepts as arguments a list of one or more file paths (e.g. ./solution.rb file1.txt file2.txt ...).
2. The program also accepts input on stdin (e.g. cat file1.txt | ./solution.rb).
3. The program outputs a list of the 100 most common three word sequences in the text, along with a count of how many times each occurred in the text. For example: `231 - i will not, 116 - i do not, 105 - there is no, 54 - i know not, 37 - i am not …`
4. The program ignores punctuation, line endings, and is case insensitive (e.g. “I love\nsandwiches.” should be treated the same as "(I LOVE SANDWICHES!!)")
5. The program is capable of processing large files and runs as fast as possible.
6. The program should be tested. Provide a test file for your solution.

You could try it against "Origin Of Species" as a test: http://www.gutenberg.org/cache/epub/2009/pg2009.txt .

## Tools Needed

1.) [Golang](https://golang.org/)
 

## Usage

You will first need a binary to run on your pc
```
go build -o relic
```

The program accepts input as a list of arguments or through stdin

```
./relic file1.txt file2.txt ...
```
or 
```
cat file1.txt | ./relic
```

If you use the test file: pg2009.txt located in this project your expected out should be:

```
Worker 1 finished scanning pg2009.txt
1.) 277 - of the same
2.) 113 - conditions of life
3.) 103 - the same species
4.) 98 - in the same
5.) 93 - species of the
.
.
.
```

## Testing

In order to run the tests in this project simply do:

```
go test -v main_test.go main.go
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[Affero General Public License v3.0](https://choosealicense.com/licenses/agpl-3.0/)