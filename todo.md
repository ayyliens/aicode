Forking support seems incomplete: it needs to refresh the backup after each FS event rerun. TODO recheck.

Instead of forbidding "holes" in file indexes, use them. When performing a run, instead of always continuing from the last file, continue from the last one before the first "hole". This would allow additional potentially useful scenarios, such as:

* Pre-create multiple "user" message files with indexes 0, 2, 4, etc.,
  as a "dialogue framework" for a sequence of expected bot responses.

* Deleting a message in the middle of a conversation would be a convenient way
  to "retry"/"redo" that part of the conversation, especially when truncation
  and/or forking is not enabled.
