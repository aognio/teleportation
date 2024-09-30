# Teleportation Changlog

## [0.1.0] - 2024-09-30

### **Added**

* Ability to create aliases for frequently accessed directories.
* Support for updating an existing alias using the `--update` flag.
* Option to delete an alias using the `--delete` flag.
* Recall and navigate to a directory using an alias, with optional integration into the shell via the `--sourced` flag.
* Persistent alias storage using SQLite for reliable management across sessions.
* Ability to list all aliases or filter them by a specific folder path using the `--list` flag.
* Option to rename an alias via the `--rename` flag.
* Rename all aliases under a specific folder path using the `--rename-folder` flag.
* Track the number of times an alias has been invoked.
* Direct access to the SQLite database through the `--sqlite` flag for manual editing.
* Informative and error messages support shell sourcing, ensuring consistent formatting when `--sourced` is used.
* Special character escaping in sourced output to avoid shell syntax errors.
