## linkwarden completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(linkwarden completion bash)

To load completions for every new session, execute once:

#### Linux:

	linkwarden completion bash > /etc/bash_completion.d/linkwarden

#### macOS:

	linkwarden completion bash > $(brew --prefix)/etc/bash_completion.d/linkwarden

You will need to start a new shell for this setup to take effect.


```
linkwarden completion bash
```

### Options

```
  -h, --help[=true|false]              help for bash
      --no-descriptions[=true|false]   disable completion descriptions
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden completion](linkwarden_completion.md)	 - Generate the autocompletion script for the specified shell

