# monkey-lang-interpreter
An interpreter for the monkey language described in ["Writing an Interpreter in Go" by Thorsten Ball](https://interpreterbook.com).

## Running
Clone/download the repo into your go workspace, navigate into the root directory of the project, and run "go run main.go"

## Testing
There are testing modules for all major objects in the monkey-lang-interpreter (ex parser, evaluator, etc). To run the test cases run the test command on the folder of interest. For example: "go test ./evaluator" will run the unit tests for the evaluator.


## Current State / Future Additions
Currently I have the lexer, and parser fully implemented with some of the evaluator done. The evaluator currently supports variables, numeric operations, and boolean operations.

I'm currently working on adding function declarations and calls to monkeylang.

Later, once I have completed the book, I plan on adding loops to monkeylang. Loops are not in the original spec of monkeylang and are not described in the book, however I view loops as a critical part of most modern programming languages and therefore should be included. 

Once the foundation on the language is built I would also like to go back and add improvements to the reliability and speed of the interpreted language.
