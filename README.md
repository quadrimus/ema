# EMA — Enhanced Markup for Authors
[![Go Reference](https://pkg.go.dev/badge/quadrimus.com/ema.svg)](https://pkg.go.dev/quadrimus.com/ema)
[![Go Report Card](https://goreportcard.com/badge/quadrimus.com/ema)](https://goreportcard.com/report/quadrimus.com/ema)

- utility and library written in pure [Go](https://go.dev/)
- `.ema` file format

## File Format
EMA file is a text-based format in UTF-8 encoding with `.ema` extension.

Each file starts with header, which is four bytes: `{EMA` (values 123, 69, 77, 65)

File content can be separated to arbitrary numbers of parts, 
with each part being one of the following:

### Text part
All bytes are interpreted as plain text.
The only special character is `{` (U+007B), which has the following meaning:
- If there is another `{` character immediately after `{`, text part
  continues and both characters are interpreted as one occurrence of `{` character.
- Otherwise, text part ends and content continues as **data part**.
  Character `{` is included to **data part**.

### Data part
All bytes are interpreted as JSON object (see [json.org](https://json.org)).
After end `}` character of object content continues as **text part**.

There is also special data part form called **command** 
where starting `{` is followed by letter `A`-`Z` or `a`-`z` 
(first command name letter).
Spaces between `{` and command name are allowed and have no meaning.
Command name match regular expression:
```regexp
^[A-Za-z][-0-9A-Za-z]*$
```
After command name, there can be optional spaces followed by:
- end `}` character of command
- JSON array followed by optional spaces and end `}` character of command
- JSON object followed by optional spaces and end `}` character of command

Spaces are all characters in Unicode category *Space Separator*.

Following commands and JSON object are equivalent:
```
{}                = {}
{bold}            = {"use": ["bold"]}
{ bold }          = {"use": ["bold"]}
{bold[1]}         = {"use": ["bold", 1]}
{ bold [1] }      = {"use": ["bold", 1]}
{bold[1, 2]}      = {"use": ["bold", 1, 2]}
{ bold [1, 2] }   = {"use": ["bold", 1, 2]}
{bold{"w": 1}}    = {"use": ["bold", {"w": 1}]}
{ bold {"w": 1} } = {"use": ["bold", {"w": 1}]}
```
