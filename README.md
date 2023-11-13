# Crontable Reader and Writer
Crontable Reader and Writer is a Go project that provides robust tools for reading, parsing, validating, and explaining cron expressions. The project is split into two main packages: reader and writer. The reader package focuses on interpreting and validating cron expressions, while the writer package converts these expressions into human-readable descriptions.

## Features
- Cron Expression Parsing: Efficiently reads and parses cron expressions.
- Validation: Validates cron expressions against syntax rules and value bounds.
- Decoding: Converts cron expression strings into a structured format.
- Human-Readable Descriptions: Translates cron expressions into clear, understandable text.

## Installation
To use the Crontable Reader and Writer in your Go project, import the packages as follows:

```
import (
    "github.com/dark-enstein/crontable/reader"
    "github.com/dark-enstein/crontable/writer"
)
```

## Usage
### Reading and Validating Cron Expressions
```
// Open a cron file and read the first line as a CronRead
cronRead, err := reader.OpenCrontableFile("path/to/cronfile")
if err != nil {
    log.Fatal(err)
}

// Validate the cron expression
isValid, err := cronRead.Validate()
if err != nil {
    log.Fatal(err)
}

// Decode the cron expression into a structured format
decoded := cronRead.Decode()
```

### Writing Human-Readable Descriptions
```
// Explain the decoded cron expression
meaning := writer.Explain(decoded)

// Output the explanation
fmt.Println(string(meaning))
```

## Dependencies
Go standard library

## Contributing
Contributions to the Crontable are welcome.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
