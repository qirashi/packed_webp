# packed_webp [RU](README.md) | [EN]
A console utility for unpacking and packing `.packed.webp` files.

## Usage

### Portable
1. Download `riff.exe` or extract the archive with the utility to a convenient location.
2. Unpack the `webp` file with `extr` from `dvpl`.
3. Drag and drop one or more files onto `riff.exe`.
   - Recursive processing of folders and files is supported.

### Console, Context Menu, and PATH Environment

#### : Console :
1. Open `cmd` or `PowerShell` in the folder with the utility.
2. Run file or folder processing:
   ```cmd
   riff.exe image.webp
   riff.exe folder\images
   ```
   You can specify multiple files and folders at once:
   ```cmd
   riff.exe file1.webp file2.packed.webp folder1 folder2
   ```

#### : PATH Environment :
To run `riff.exe` from any folder:
1. Add your folder to the **Environment Variables**, for example `C:\Tools\packed_webp`, to the PATH:
   - Win+R → `sysdm.cpl` → **Advanced** tab → **Environment Variables**.
   - Find the `PATH` variable, click **Edit**, and add:
   ```
   ;C:\Tools\packed_webp
   ```
2. Now you can simply run:
   ```cmd
   riff.exe file.webp
   riff.exe folder
   ```

#### : Context Menu (Windows) :
To run processing directly from the right-click menu:
1. Copy `riff.exe` to a convenient location, e.g., `C:\Tools\packed_webp`.
2. Create a text file `packed_webp.reg` with the following content:
   ```reg
   Windows Registry Editor Version 5.00

   [HKEY_CLASSES_ROOT\*\shell\packed_webp]
   @="Process packed_webp"

   [HKEY_CLASSES_ROOT\*\shell\packed_webp\command]
   @="\"C:\\Tools\\packed_webp\\riff.exe\" \"%1\""
   ```
3. Save the file and double-click `packed_webp.reg` to apply the changes.
4. Now the **Process packed_webp** option will appear in the file context menu.

---
