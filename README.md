# Description

Creates a file based upon a programming language template. Essentially touch but with added powers. Custom scripting
is available using the Lua programming language.

## Example Terminal Output

![Example Terminal Output](res/touchy-example.png)

## Commands

All commands take a programming language and the template name as arguments for the most part. In most cases the template name can be omitted.

## Create

Creates a new file based on a template. The example below creates a file based on the default 'go' language template.

```text
touchy create go
```

You can also provide a template name after the language name to create a specific template:

```text
touchy create go main
```

## List

List all available programming languages and/or templates of a programming language. The example below lists all programming languages.

```text
touchy list all
```

To get a list of a programming languages templates use the following command example.

```text
touchy list go
```

To get a list of scripts use the following command example

```text
touchy list scripts
```

## Show

Displays the colorized contents of a template to the terminal.
Example that will show the contents of the default go template.

```text
touchy show go
```

You can also provide a template name after the language name to show a specific template.

```text
touchy show go main
```

## Run

Runs a script. The example below runs the default test script.

```text
touchy run test
```

Again you can get a list of scripts using the following command:

```text
touchy list scripts
```

*Note that the API is not yet stablized and will be documented once it is.*
