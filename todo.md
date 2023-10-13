Support response streaming.

Forking support seems incomplete: it needs to refresh the backup after each FS event rerun. TODO recheck.

Support forking in non-watch mode from specific index.

Merge `--trunc` and `--trunc-after` into one, where watch mode supports boolean or index, and non-watch mode supports only index.

Instead of forbidding "holes" in file indexes, use them. When performing a run, instead of always continuing from the last file, continue from the last one before the first "hole". This would allow additional potentially useful scenarios, such as:

* Pre-create multiple "user" message files with indexes 0, 2, 4, etc.,
  as a "dialogue framework" for a sequence of expected bot responses.

* Deleting a message in the middle of a conversation would be a convenient way
  to "retry"/"redo" that part of the conversation, especially when truncation
  and/or forking is not enabled.

Every operation involving the bot could be run redundantly, concurrenly, with configurable N concurrency. Then compare the outputs. The bot is fuzzy and the outputs for the exact same inputs may differ between runs.

Support working in multiple directories at once. The framework could watch 1 or N ancestor directories, containing N conversation directories as children or descendants. This allow the user to concurrently work on multiple different features / prompts while waiting for bot responses, which can take minutes. The user would make changes in one directory to launch a request changes, make changes in another directory to launch another concurrent request, and so on.

Support optional Go files as "plugins" in conversation directories. When launching the framework, instead of only treating directories as data, we could look for Go files and add them on the command line as arguments to the `go` tool. The files would import the framework and "plug" into it, configuring the framework for their directories, like configuring and registering supported functions.

Automatic prompt development. Currently, prompt development is the bottleneck. Generating a useful prompt for any given task can take many attempts, and each attempt can take several minutes, mostly waiting for bot responses. We could automate this by running a loop similar to the following:

* Step 0. Write prompt, test, correction phrase. Loop:
* Step 1. Provide prompt and inputs to the bot.
* Step 2. Run a test, which may return a boolean or error text.
* Step 3. If the test fails, provide the latest prompt, the inputs, the outputs, the error text if any, and tell the bot to generate a different prompt for the same task. Goto 1 with the new prompt.

↑ For any given task, the original human developer may have to write the original prompt, the test, and something like a correction phrase which will be used in step 3. The test may be anything that takes the inputs and outputs and generates error text. The test may be regular code that runs locally. The test may involve sending the inputs and outputs to the bot, asking to compare them.
