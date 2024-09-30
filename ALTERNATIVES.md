Alternatives to Teleportation
=============================

Teleportation is a command-line utility for managing and navigating folders and their associated aliases efficiently. Here are some existing alternatives to Teleportation, which provide similar functionality for navigating the filesystem using shortcuts:

1. [zoxide](https://github.com/ajeetdsouza/zoxide)
--------------------------------------------------

**Summary**: `zoxide` is a smarter `cd` command that learns your habits. It uses a ranking algorithm to prioritize frequently visited directories and allows you to jump between them quickly.

**Features**:

* Machine-learning-like ranking based on frequency and recency.
* Supports all major shells (bash, zsh, fish, etc.).
* Commands are simple: `z <partial-dir-name>` to teleport to a directory.

**Pros**:

* Fast and lightweight, implemented in Rust.
* No configuration required; it learns your habits out of the box.

**Cons**:

* Limited customization for aliases, as it focuses more on frequency-based navigation.

* * *

2. [autojump](https://github.com/wting/autojump)
------------------------------------------------

**Summary**: `autojump` is a command-line tool that also helps users navigate directories by learning which folders they frequent the most.

**Features**:

* Tracks which directories you visit the most, allowing you to jump to them quickly.
* Integrates well with most shells.

**Pros**:

* Automatically learns your directory preferences over time.
* Can be customized with `autojump`'s internal ranking.

**Cons**:

* Requires some training time before becoming truly efficient.
* May be slower in environments with a very large number of directories.

* * *

3. [z](https://github.com/rupa/z)
---------------------------------

**Summary**: `z` is a utility that helps you quickly jump between directories you've frequently or recently accessed, based on a scoring mechanism.

**Features**:

* Lightweight script with no dependencies.
* Jumps to directories based on both frequency and recency.

**Pros**:

* Simple to use and configure.
* Works well for users who prefer minimal dependencies.

**Cons**:

* Limited support for custom aliases.
* Script can be slower on very large systems, as it uses a basic scoring system.

There is also [Zsh](https://www.zsh.org/) version called [Zsh-z](https://github.com/agkozak/zsh-z) too.

* * *

Comparison with Teleportation
-----------------------------

**Teleportation**: Unlike the above tools, Teleportation focuses on explicitly defined aliases, allowing users to manage directory shortcuts more efficiently, without relying on a learning algorithm. Users can easily create, modify, and delete these aliases, making it highly customizable.

**Key Differences**:

1. **Alias Management**: Teleportation is designed around user-defined aliases, which allows for more explicit control compared to the ranking-based navigation of other tools like `zoxide`, `autojump`, and `z`.
2. **Customization**: Teleportation’s alias system offers greater flexibility for organizing shortcuts in ways that are meaningful to the user, rather than relying on implicit behavior.
3. **Use Case**: While `zoxide`, `autojump`, and `z` are great for users who prefer their tool to "learn" their habits, Teleportation is perfect for those who want a reliable set of static shortcuts without any automated ranking.

Other popular techniques for jumping to other directories with shortcuts
------------------------------------------------------------------------

### 2. **Aliases**

You can use shell aliases to create shortcuts for directories you use often.

```sh
alias proj="cd /home/user/documents/projects"
```

This way, typing `proj` will take you to that directory.

### 3. **Bookmark-like Scripts**

You can write your own scripts to create bookmarks for directories, often called `marks`. For instance:

* Create a file to store directory shortcuts (`~/.marks`).
* Write functions to save and load those bookmarks.

Example:

```sh
mark() {
  echo "export $1=\"$(pwd)\"" >> ~/.marks
}

go() {
  source ~/.marks
  cd "$(eval echo \$$1)"
}
```

You can then run `mark proj` to save the current directory as `proj` and use `go proj` to return to it.

### 4. **`pushd` and `popd`**

These are built-in shell commands that allow you to maintain a directory stack.

* `pushd <directory>` pushes the directory onto the stack and changes to it.
* `popd` pops the top directory off the stack and changes back to it.

### 5. **CDPATH Environment Variable**

`CDPATH` is a shell feature that makes it easy to change directories to common locations without typing the full path.

For example:

```sh
export CDPATH=.:/home/user/documents:/var/www
```

Now you can just type `cd projects` if `projects` is inside `/home/user/documents`, and it will take you there.

These tools and techniques make directory navigation much more convenient in UNIX-like environments.



