# PLAIN Standard Library Specification ("Batteries Included")

This document outlines the core module layout, implemented runtime namespaces, and candidate function sets for the PLAIN standard library, targeting the desktop application and business utility niches historically occupied by tools like Visual Basic, Delphi, and Xojo.

> **NOTE:** Modules marked as **Implemented** are active in the current runtime. All other modules represent planned design specifications.

---

## Module Architecture

Standard library modules are exposed as global `Struct` namespaces pre-registered at compile time and bound to low-index global slots in the VM (`System` at slot `0`, `Console` at slot `1`, etc.). Module functions are implemented as native Go `ValBuiltin` callbacks invoked via standard dot-notation property access.

---

### Canvas *(Reserved / Future)*

* **Purpose**: Simple 2D graphics canvas for lightweight retro games and visual rendering.
* **Functions**:
  * `Box(x1, y1, x2, y2, color)`
  * `Circle(x, y, radius, color)`
  * `Create(width, height)`
  * `Fill(color)`
  * `Line(x1, y1, x2, y2, color)`
  * `Point(x, y, color)`

### Console *(Implemented — Global Slot 1)*

* **Purpose**: Standard terminal input, output, screen formatting, and ANSI escape code control.
* **Implemented Functions**:
  * `Cls()` — Clear the terminal screen and reset cursor position.
  * `Input(prompt)` — Print prompt text and read a line of string input from `stdin`.
  * `Print(text)` — Output text to `stdout`.
  * `PrintLine(text)` — Output text followed by a newline to `stdout`.
* **Planned Functions**:
  * `CursorTo(row, col)`
  * `ResetColor()`
  * `SetColor(foreground, background)`

### CSV *(Planned)*

* **Purpose**: Comma-separated values parsing and tabular data generation.
* **Functions**:
  * `Parse(csvString)`
  * `ReadRow(fileStream)`
  * `Stringify(tableData)`
  * `WriteRow(fileStream, rowData)`

### Database *(Planned)*

* **Purpose**: Embedded SQLite database operations for lightweight local persistence.
* **Functions**:
  * `BeginTx(db)`
  * `Close(db)`
  * `Commit(tx)` / `Rollback(tx)`
  * `Execute(db, sql, params...)`
  * `Open(path)`
  * `Query(db, sql, params...)`

### Debugging *(Planned)*

* **Purpose**: Runtime inspection and debugging support.
* **Functions**:
  * *TBD*

### Decimal *(Implemented — Global Slot 9)*

* **Purpose**: Rounding and fixed-place formatting helpers for PLAIN `Numeric` values. Values remain IEEE 754 `float64`; this module does not provide fixed-point or arbitrary-precision arithmetic.
* **Implemented Functions**:
  * `Format(value, places)` — Return a string with exactly `places` digits after the decimal point (`0` through `1000`).
  * `Round(value, places)` — Round a numeric value to the requested decimal place; negative places round to tens, hundreds, and so on (`-308` through `308`).

### Email *(Candidate / Future)*

* **Purpose**: Automated SMTP messaging and report delivery for business alerts and workflows.
* **Functions**:
  * `AttachFile(message, filePath)`
  * `ConfigureSMTP(host, port, username, password)`
  * `CreateMessage(to, subject, body)`
  * `Send(message)`

### Environment *(Planned)*

* **Purpose**: Platform-specific system inspection, path resolution, and variable access.
* **Functions**:
  * `GetArchitecture()`
  * `GetEnv(name)` / `SetEnv(name, value)`
  * `GetPlatform()`
  * `GetTempDir()` / `GetHomeDir()`
  * `IsBigEndian()`

### File *(Implemented — Global Slot 7)*

* **Purpose**: Synchronous whole-file text and byte input/output.
* **Implemented Functions**:
  * `AppendText(path, content)` — Append text to a file, creating it if needed.
  * `Exists(path)` — Return whether the path exists; filesystem errors other than “not found” are reported.
  * `ReadBytes(path)` — Return file bytes as a list of numeric values from `0` to `255`.
  * `ReadText(path)` — Read a file as a string.
  * `WriteBytes(path, byteValues)` — Write a list of whole-number byte values from `0` to `255`.
  * `WriteText(path, content)` — Create or replace a file with text.

### Forms *(Reserved / Future)*

* **Purpose**: Platform-dependent native GUI windowing and control management.
* **Functions**:
  * `AddButton(window, label, onClick)`
  * `AddTextBox(window, initialText)`
  * `CreateWindow(title, width, height)`
  * `ShowModal(window)`

### Image *(Planned)*

* **Purpose**: Loading, saving, and manipulating raster image assets.
* **Functions**:
  * `Crop(image, x, y, width, height)`
  * `GetPixel(image, x, y)`
  * `Load(path)`
  * `Resize(image, width, height)`
  * `Save(image, path)`
  * `SetPixel(image, x, y, color)`

### Input *(Planned)*

* **Purpose**: Low-level hardware polling for keyboards, mice, and joysticks.
* **Functions**:
  * `GetJoystickAxis(joyId, axisId)`
  * `GetMouseX()` / `GetMouseY()`
  * `IsKeyDown(keyCode)`
  * `IsMouseButtonDown(buttonCode)`

### JSON *(Planned)*

* **Purpose**: JSON parsing, serialization, and structure traversal.
* **Functions**:
  * `GetPath(object, pathExpression)`
  * `Parse(jsonString)`
  * `Stringify(value)`

### List *(Implemented — Global Slot 2)*

* **Purpose**: Basic collection operations for PLAIN lists.
* **Implemented Functions**:
  * `Contains(list, value)`
  * `Len(list)`
  * `Append(list, value)` — Return a list with `value` appended.
  * `Reverse(list)`
* **Indexing**: List elements use 1-based indices, from `1` through `List.Len(list)`.

### Locale *(Candidate / Future)*

* **Purpose**: Regional formatting for dates, numbers, thousands separators, and currency symbols.
* **Functions**:
  * `FormatCurrency(amount, localeCode)`
  * `FormatDate(timestamp, localeCode)`
  * `FormatNumber(number, localeCode)`
  * `GetSystemLocale()`

### Log *(Candidate / Future)*

* **Purpose**: Structured application logging and audit trails with severity levels and file redirection.
* **Functions**:
  * `Debug(message)`
  * `Error(message)`
  * `Info(message)`
  * `SetLevel(level)`
  * `SetOutputFile(path)`
  * `Warn(message)`

### Math *(Implemented — Global Slot 5)*

* **Purpose**: Common numeric operations.
* **Implemented Functions**:
  * `Abs(n)`, `Ceil(n)`, `Floor(n)`, `Round(n [, decimals])`
  * `Max(a, b)`, `Min(a, b)`, `Pow(base, exp)`, `Sqrt(n)`

### Network *(Planned)*

* **Purpose**: Lightweight HTTP client capabilities for web requests and API integration.
* **Functions**:
  * `Delete(url)`
  * `Get(url)`
  * `Post(url, body)`
  * `Put(url, body)`
  * `SetHeader(key, value)`

### Random *(Implemented — Global Slot 6)*

* **Purpose**: Pseudo-random number generation and list selection.
* **Implemented Functions**:
  * `Chance(percentage)` — Return a Boolean with the supplied probability from `0` to `100`.
  * `Choice(list)` — Return one list element, or `None` for an empty list.
  * `Float()` — Return a number in `[0, 1)`.
  * `Int(min, max)` — Return an integer in the inclusive range; reversed bounds are accepted.
  * `Seed(number)` — Reset the shared generator using a finite signed 64-bit integer seed.
  * `Shuffle(list)` — Return a shuffled copy of the list.

### Sound *(Planned)*

* **Purpose**: Audio playback for sound effects and basic background tracks.
* **Functions**:
  * `LoadSound(path)`
  * `PlaySound(sound)`
  * `SetVolume(sound, level)`
  * `StopSound(sound)`

### String *(Implemented — Global Slot 3)*

* **Purpose**: Basic Unicode-aware text manipulation.
* **Implemented Functions**:
  * `Join(list, delimiter)` — Join string representations of list elements.
  * `Left(str, count)`, `Right(str, count)` — Return up to `count` Unicode code points from either end.
  * `Len(str)` — Count Unicode code points.
  * `Lower(str)`, `Upper(str)`
  * `Mid(str, start [, length])` — Extract from a 1-based start position.
  * `Replace(str, target, replacement)` — Replace all occurrences.
  * `Split(str, delimiter)` — Split into a list of strings.

### Struct *(Implemented — Global Slot 4)*

* **Purpose**: Basic reflection over struct values.
* **Implemented Functions**:
  * `HasKey(structInstance, keyName)` — Return whether the string key exists.
  * `Keys(structInstance)` — Return the field names in lexicographic order.

### System *(Implemented — Global Slot 0)*

* **Purpose**: Runtime error interception, process termination, shell spawning, and OS-level control.
* **Implemented Functions**:
  * `Exit(code)` — Stop script execution; the command-line runner exits with the supplied numeric status code. Exit requests are not passed to the runtime error handler.
  * `SetErrorHandler(handlerFn)` — Register a global callback function or closure to intercept VM runtime faults.
* **Planned Functions**:
  * `Exec(command, args...)`
  * `GetPID()`
  * `Sleep(milliseconds)`
  * `Spawn(command, args...)`

### Testing *(Planned)*

* **Purpose**: Software testing and assertion support.
* **Functions**:
  * *TBD*

### Time *(Implemented — Global Slot 8)*

* **Purpose**: Millisecond timestamps and basic delays/formatting.
* **Implemented Functions**:
  * `Format(timestamp)` — Format a Unix-millisecond timestamp as `YYYY-MM-DD HH:MM:SS`.
  * `Now()` — Return the current Unix timestamp in milliseconds.
  * `Sleep(milliseconds)` — Suspend execution for the requested non-negative duration.

### XML *(Planned)*

* **Purpose**: Structured XML document parsing and querying.
* **Functions**:
  * `Parse(xmlString)`
  * `Query(node, xpath)`
  * `ToString(node)`
