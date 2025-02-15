# simple-markov-go

**simple-markov** is a command-line tool that reads text from either a file or standard input, builds a character-level Markov model of a specified order, and generates random text of a specified length. You can also provide a **starter** text using the `-starter` flag, which is re-output at the beginning of the generated text.

## Features

- **Character-level** Markov chain.
- **Configurable order** (`k`).
- **Configurable output length** (`l`).
- **Optional starter text** (`-starter`).
- Reads from **stdin** or from an **input file** (`-i`).
- **Fast** and memory-efficient for typical text processing tasks.

## Installation

1. **Clone** this repository:
   ```bash
   git clone https://github.com/yourusername/simple-markov.git
   cd simple-markov
   ```
2. **Build** the binary:
   ```bash
   go build -o simple-markov .
   ```
   This will create an executable named `simple-markov`.

## Usage

You can provide input via standard input:

```bash
cat text.txt | ./simple-markov -k 1 -l 300
```

Or specify an input file with `-i`:

```bash
./simple-markov -k 1 -l 300 -i text.txt
```

Available flags:

- `-k int` : The order (or "look-back") of the Markov chain. Default is `1`.
- `-l int` : The length (in characters) of output to generate. Default is `100`.
- `-i string` : The file path of the input text. If omitted, the program reads from **stdin**.
- `-seed int` : If provided, uses a fixed seed for reproducible outputs. If omitted, uses the current time for the seed.
- `-starter string` : A string that will be placed at the beginning of the generated text, and then the Markov model continues from the last `k` characters of this starter text.

## Example

```bash
# Build a Markov model of order 2, start with "Hello ", and generate 50 total characters
cat your_text_corpus.txt | ./simple-markov -k 2 -l 50 -starter "Hello "
```

This will produce output that begins with `"Hello "` and then continues up to 50 total characters.

---

## Contributing

Feel free to open a pull request if you have any suggestions or improvements.

---

## License

This project is provided under the MIT License.