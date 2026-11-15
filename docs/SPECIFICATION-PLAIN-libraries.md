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

### Decimal *(Candidate / Future)*

* **Purpose**: Fixed-point arithmetic and rounding modes for accurate financial, accounting, and currency calculations without floating-point discrepancies.
* **Functions**:
  * `Add(a, b)`, `Sub(a, b)`, `Mul(a, b)`, `Div(a, b)`
  * `FormatCurrency(amount, pictureClause)`
  * `FromString(numericString)`
  * `RoundBankers(amount, decimals)`

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

### File *(Planned)*

* **Purpose**: Synchronous and stream-based file system input/output.
* **Functions**:
  * `Append(path, data)`
  * `Copy(src, dest)`
  * `Delete(path)`
  * `Exists(path)`
  * `Open(path, mode)`
  * `Read(file)` / `ReadAll(path)`
  * `Write(file, data)` / `WriteAll(path, data)`

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

### Math *(Planned)*

* **Purpose**: Business, programmer, and scientific mathematical operations.
* **Functions**:
  * `Abs(n)`, `Floor(n)`, `Ceil(n)`, `Round(n, decimals)`
  * `Log(n)`, `Exp(n)`, `Pow(base, exp)`, `Sqrt(n)`
  * `Min(a, b)`, `Max(a, b)`
  * `Sin(n)`, `Cos(n)`, `Tan(n)`, `Asin(n)`, `Acos(n)`, `Atan(n)`

### Network *(Planned)*

* **Purpose**: Lightweight HTTP client capabilities for web requests and API integration.
* **Functions**:
  * `Delete(url)`
  * `Get(url)`
  * `Post(url, body)`
  * `Put(url, body)`
  * `SetHeader(key, value)`

### Sound *(Planned)*

* **Purpose**: Audio playback for sound effects and basic background tracks.
* **Functions**:
  * `LoadSound(path)`
  * `PlaySound(sound)`
  * `SetVolume(sound, level)`
  * `StopSound(sound)`

### String *(Planned)*

* **Purpose**: Advanced text manipulation, substring extraction, pattern matching, and encoding utilities.
* **Functions**:
  * `Base64Encode(str)` / `Base64Decode(str)`
  * `ExtractRegex(str, pattern)`
  * `Join(list, delimiter)`
  * `Left(str, count)`
  * `Len(str)`
  * `Lower(str)` / `Upper(str)`
  * `MatchRegex(str, pattern)`
  * `Mid(str, start, length)`
  * `Replace(str, target, replacement)`
  * `ReplaceRegex(str, pattern, replacement)`
  * `Right(str, count)`
  * `Split(str, delimiter)`

### Struct *(Planned)*

* **Purpose**: Structural reflection, key enumeration, and layout serialization.
* **Functions**:
  * `FromJson(jsonString, structTypeName)`
  * `HasKey(structInstance, keyName)`
  * `Keys(structInstance)`
  * `ToJson(structInstance)`

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

### Time *(Planned)*

* **Purpose**: Temporal calculations, formatting, and date-time manipulation.
* **Functions**:
  * `AddDays(timestamp, days)`
  * `Diff(timestampA, timestampB)`
  * `Format(timestamp, layout)`
  * `Now()`
  * `Parse(dateString, layout)`
  * `Year(t)`, `Month(t)`, `Day(t)`, `Hour(t)`, `Minute(t)`, `Second(t)`

### XML *(Planned)*

* **Purpose**: Structured XML document parsing and querying.
* **Functions**:
  * `Parse(xmlString)`
  * `Query(node, xpath)`
  * `ToString(node)`