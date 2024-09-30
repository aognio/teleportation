# Teleportation (e.g., `tlp2-bin`)

`Teleportation` is a command-line utility for managing and navigating folders and their associated aliases efficiently. It is designed to save time by allowing users to quickly "teleport" to directories using predefined aliases.

## Purpose

The purpose of `tlp2-bin` is to provide a fast and easy way to manage and navigate to directories using short aliases. This can be especially useful for users who frequently switch between deep or complex directory structures. By creating an alias for a folder, you can quickly jump to that location without typing the entire path.

## Features

- **Store Alias for Folders**: Create an alias for any folder, allowing rapid access to that location.
- **Update Existing Alias**: Modify the path associated with an existing alias using the `--update` flag.
- **Delete Alias**: Remove an alias from the database using the `--delete` flag.
- **Recall Folder by Alias**: Use an alias to print or navigate directly to a folder.
- **Shell Integration with `--sourced` Flag**: Output commands formatted for direct sourcing by your shell to facilitate navigation.
- **Persistent Storage with SQLite**: Alias information is stored persistently using an SQLite database.
- **Open SQLite Database Manually**: Use the `--sqlite` flag to access the SQLite database directly.
- **List Aliases**: List all saved aliases or list aliases for a specific folder path using the `--list` flag.
- **Rename Alias**: Rename an existing alias with the `--rename` flag.
- **Rename Folder Path**: Update all aliases under a specific folder path using the `--rename-folder` flag.
- **Invocation Count Tracking**: Track the number of times an alias has been invoked for usage insights.
- **Special Character Escaping**: Automatically escape special characters in output to avoid shell syntax errors.

## Usage

Create an alias like this:

```bash
tlp2 code/freelancing/projects petstore_website /home/johndoe/freelancing/projects/sarah/petstore/website
```

Now you can quickly jump to that folder just by typing this command:

```bash
tlp2 petstore_website
```

If there were two `petstore_website` aliases, you will have to specify the folder:

```bash
tlp2 code/freelancing/projects petstore_website
```

And that's it, for the most part. 

The actual command binary is wrapped in a shell function so an actual "cd /path/to/folder" command is part of the output but your shell executes the command and the directory change takes places.

There are, of course, many other uses cases related to mantaining this aliases collections that are detailed in the next section.

Everything is stored in a SQLite database located at:

```
$HOME/.config/teleportation/teleportation.sqlite3
```

If you have the `sqlite3` command installed you can manipulate the records directly using SQL by issuing this command:

```bash
tlp2 --sqlite
```

### Use cases

```bash
tlp2 code/projects myapp /home/user/code/myapp
```
#### 1. Create alias:
Expected output:
```
Alias 'myapp' created.
```

#### 2. Update alias:
```bash
tlp2 --update code/projects myapp /home/user/newpath/myapp
```
Expected output:
```
Alias 'myapp' updated.
```

#### 3. Delete alias:
```bash
tlp2 --delete code/projects myapp
```
Expected output:
```
Alias 'myapp' under folder 'code/projects' has been deleted.
```

#### 4. Recall alias:
```bash
tlp2 code/projects myapp
```
Effect:
```
The parent shell changes directories to `/home/user/code/myapp`.
```

#### 5. Search for alias across all folders:
```bash
tlp2 myapp
```
Effect:
```
The parent shell changes directories to the saved path associated with the alias.
```

#### 6. List all aliases:
```bash
tlp2 --list
```
Expected output:
```
5 aliases created under 3 different folders:
alias1 folder1 /path/to/folder1
alias2 folder2 /path/to/folder2
```

#### 7. List aliases under specific folder:
```bash
tlp2 --list code/projects
```
Expected output:
```
2 aliases found under 1 different folder:
alias1 code/projects /home/user/code/projects/alias1
alias2 code/projects /home/user/code/projects/alias2
```

#### 8. Rename folder:
```bash
tlp2 --rename-folder old/folder new/folder
```
Expected output:
```
Folder 'old/folder' has been renamed to 'new/folder'.
```

### Alternative execution flows for the use cases (e.g. errors)

#### 1. Create alias without update flag (alias exists):
```bash
tlp2 code/projects myapp /home/user/code/anotherpath
```
Expected output:
```
Error: Alias 'myapp' under folder 'code/projects' already exists. Use --update to modify it.
```

#### 2. Delete alias that does not exist:
```bash
tlp2 --delete code/projects nonexistent
```
Expected output:
```
Error: Alias 'nonexistent' under folder 'code/projects' does not exist.
```

#### 3. Recall non-existent alias:
```bash
tlp2 code/projects nonexistent
```
Expected output:
```
Error: Alias 'nonexistent' under folder 'code/projects' does not exist.
```

#### 4. List aliases under non-existent folder:
```bash
tlp2 --list nonexistent/folder
```
Expected output:
```
No folders match the 'nonexistent/folder' value
```

#### 5. Rename non-existent folder:
```bash
tlp2 --rename-folder nonexistent/folder new/folder
```
Expected output:
```
Error: Folder 'nonexistent/folder' does not exist.
```

#### 6. Invalid usage (incorrect parameters):
```bash
tlp2 code/projects
```
Expected output:
```
Usage for storing: tlp2 <folder_path> <alias> <absolute_path>
Usage for updating: tlp2 --update <folder_path> <alias> <absolute_path>
Usage for recalling: tlp2 <folder_path> <alias>
Usage for searching: tlp2 <alias>
Usage for listing all aliases: tlp2 --list [<partial_folder>]
Usage for deleting: tlp2 --delete <folder_path> <alias>
Usage for renaming a folder: tlp2 --rename-folder <old_folder_path> <new_folder_path>
```


## How to Build

### Unix/Linux (Bash)

1. Clone the repository and navigate to the project root:

   ```sh
   git clone <repository_url>
   cd tlp2-bin
   ```

2. Run the build script:
    
   ```sh
   ./build.sh
   ```
    
3. This will create an executable named `tlp2-bin`.
    

### Windows (CMD)

1. Clone the repository and navigate to the project root:
    
   ```cmd
   git clone <repository_url>
   cd tlp2-bin
   ```
    
2. Run the CMD build script:
    
   ```cmd
   build.cmd
   ```
    
3. This will create an executable named `tlp2-bin.exe`.
    

### Windows (PowerShell)

1. Clone the repository and navigate to the project root:
    
   ```powershell
   git clone <repository_url>
   cd tlp2-bin
   ```
    
2. Run the PowerShell build script:
    
   ```powershell
   ./build.ps1
   ```
    
3. This will create an executable named `tlp2-bin.exe`.
    
## Installation

To install `tlp2-bin`, run the following command:

```sh
go install github.com/aognio/teleportation/cmd/tlp2-bin@latest
```

This command will download the latest version of the source code, compile it, and place the executable in your Go binaries directory. By default, this directory is `$GOPATH/bin` or `$HOME/go/bin` if `$GOPATH` is not set. You need to make sure that the `$GOPATH/bin` directory is included in your system's `PATH` environment variable so that you can run `tlp2-bin` from anywhere.

## Adding `$GOPATH/bin` to Your PATH

To add the Go binaries directory to your `PATH`, follow the instructions for your platform:

### Unix/Linux

1. Open your shell profile file (e.g., `.bashrc`, `.bash_profile`, or `.zshrc`) in a text editor:
    
   ```sh
   nano ~/.bashrc
   ```
    
2. Add the following line to include `$GOPATH/bin` in your `PATH`:
    
   ```sh
   export PATH=$PATH:$HOME/go/bin
   ```
    
3. Save the file and apply the changes:
    
   ```sh
   source ~/.bashrc
   ```
    
### Windows

1. Open the Start menu and search for "Environment Variables."
2. Click on "Edit the system environment variables."
3. In the System Properties window, click the "Environment Variables" button.
4. In the "System variables" section, select "Path" and click "Edit."
5. Add a new entry for `%USERPROFILE%\go\bin` or `%GOPATH%\bin` if you have set a custom Go path.

## Wrapping `tlp2-bin` in a Shell Alias

To make `tlp2-bin` more convenient to use, you can wrap it in a shell alias. This allows you to call `tlp2-bin` and pass all parameters to the wrapped command easily.

### Bash

Add the following alias to your `.bashrc` or `.bash_profile`:

```sh
alias tlp2='. tlp2-bin "$@"'
```

After adding this alias, make sure to reload your shell configuration:

```sh
source ~/.bashrc
```

### Zsh

Add the following alias to your `.zshrc`:

```sh
alias tlp2='. tlp2-bin "$@"'
```

Then, reload your Zsh configuration:

```sh
source ~/.zshrc
```

### Windows Command Prompt

For Windows Command Prompt, you can create a batch script to replicate the alias functionality:

1. Create a file named `tlp2.cmd` in a folder that is in your `PATH` with the following content:
    
   ```cmd
   @echo off
   tlp2-bin %*
   ```

### PowerShell

To create a PowerShell alias that passes all parameters to `tlp2-bin`, add the following to your profile script (`$PROFILE`):

```powershell
function tlp2 {
    & "Path\To\tlp2-bin.exe" @Args
}
```

After saving the profile script, reload it:

```powershell
. $PROFILE
```

### Running the tests

To verify the functionality of `tlp2`, a set of unit tests has been included. These tests ensure that all the core functionalities such as creating, updating, deleting, recalling, and listing aliases are working correctly. Follow the steps below to run the tests:

1. **Ensure dependencies are installed**: 
   - Make sure you have Go installed and properly set up. You can verify this with:
     ```sh
     go version
     ```
   - If `go` is not installed, follow the installation instructions for your platform on the [official Go website](https://golang.org/doc/install).
   - I particularly like managing my Golang installation using the [asdf](https://asdf-vm.com/) version manager.

2. **Navigate to the project directory**:
   - Go to the root directory of the project:
     ```sh
     cd /path/to/teleportation
     ```

3. **Run tests**:
   - Use the `go test` command to run all the test cases:
     ```sh
     go test ./cmd/tlp2-bin/
     ```
   - This command will run the unit tests defined in `tlp2_test.go` and output the results, indicating which tests passed or failed.

4. **Expected output**:
   - You should see an output summary that shows all tests passing if everything is set up correctly:
     ```
     ok  	github.com/aognio/teleportation/cmd/tlp2-bin	<time_taken>	s
     ```
   - If any test fails, you will see detailed information about what went wrong, making it easier to debug.

5. **Test coverage** (Optional):
   - To see the test coverage, you can use:
     ```sh
     go test ./cmd/tlp2-bin/ -cover
     ```
   - This will provide a coverage percentage, giving you an idea of how much of the code is covered by the tests.

## Features Summary

- **Alias Management**:
  - Create, update, delete, and rename aliases for directory paths.
- **Folder Path Management**:
  - Rename all aliases under a specific folder path.
- **Navigation**:
  - Recall and navigate to a folder by using its alias, with direct shell integration using `--sourced`.
- **List and Search**:
  - List all saved aliases or filter them by folder.
  - Search for aliases across all folders.
- **Persistent Alias Storage**:
  - Backed by SQLite, providing reliable storage of aliases across sessions.
- **Direct SQLite Access**:
  - Open the SQLite database directly with `--sqlite` for manual modifications.
- **Shell Integration**:
  - Use the `--sourced` flag to format output for direct navigation in the shell.
  - Special character escaping in output to ensure shell commands run smoothly.

## Contributing

If you'd like to contribute to `tlp2-bin`, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
