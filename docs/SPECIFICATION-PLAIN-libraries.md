# PLAIN Standard Library Proposal ("Batteries Included")

This proposal outlines the core module layout and candidate function sets for the PLAIN standard library, targeting the desktop application and business utility niches historically occupied by tools like Visual Basic, Delphi, and Xojo.

> **THIS IS PRELIMINARY DOCUMENTATION. IT IS SUBJECT TO REVISION**

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

### Console

* **Purpose**: Standard terminal input, output, screen formatting, and ANSI escape code control.
* **Functions**:
* `Cls()` (Clear screen)
* `CursorTo(row, col)`
* `Input(prompt)`
* `Print(text)` / `PrintLine(text)`
* `ResetColor()`
* `SetColor(foreground, background)`

### CSV

* **Purpose**: Comma-separated values parsing and tabular data generation.
* **Functions**:
* `Parse(csvString)`
* `ReadRow(fileStream)`
* `Stringify(tableData)`
* `WriteRow(fileStream, rowData)`

### Database

* **Purpose**: Embedded SQLite database operations for lightweight local persistence.
* **Functions**:
* `BeginTx(db)`
* `Close(db)`
* `Commit(tx)` / `Rollback(tx)`
* `Execute(db, sql, params...)`
* `Open(path)`
* `Query(db, sql, params...)`

### Environment

* **Purpose**: Platform-specific system inspection, path resolution, and variable access.
* **Functions**:
* `GetArchitecture()`
* `GetEnv(name)` / `SetEnv(name, value)`
* `GetPlatform()`
* `GetTempDir()` / `GetHomeDir()`
* `IsBigEndian()`

### File

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

### Image

* **Purpose**: Loading, saving, and manipulating raster image assets.
* **Functions**:
* `Crop(image, x, y, width, height)`
* `GetPixel(image, x, y)`
* `Load(path)`
* `Resize(image, width, height)`
* `Save(image, path)`
* `SetPixel(image, x, y, color)`

### Input

* **Purpose**: Low-level hardware polling for keyboards, mice, and joysticks.
* **Functions**:
* `GetJoystickAxis(joyId, axisId)`
* `GetMouseX()` / `GetMouseY()`
* `IsKeyDown(keyCode)`
* `IsMouseButtonDown(buttonCode)`

### JSON

* **Purpose**: JSON parsing, serialization, and structure traversal.
* **Functions**:
* `GetPath(object, pathExpression)`
* `Parse(jsonString)`
* `Stringify(value)`

### List

* **Purpose**: High-order collection utilities for searching, sorting, and transforming lists.
* **Functions**:
* `Contains(list, value)`
* `Filter(list, predicateFn)`
* `Find(list, predicateFn)`
* `Map(list, transformFn)`
* `Reverse(list)`
* `Sort(list, comparatorFn)`

### Math

* **Purpose**: Business, programmer, and scientific mathematical operations.
* **Functions**:
* `Abs(n)`, `Floor(n)`, `Ceil(n)`, `Round(n, decimals)`
* `Log(n)`, `Exp(n)`, `Pow(base, exp)`, `Sqrt(n)`
* `Min(a, b)`, `Max(a, b)`
* `Sin(n)`, `Cos(n)`, `Tan(n)`, `Asin(n)`, `Acos(n)`, `Atan(n)`

### Network

* **Purpose**: Lightweight HTTP client capabilities for web requests and API integration.
* **Functions**:
* `Delete(url)`
* `Get(url)`
* `Post(url, body)`
* `Put(url, body)`
* `SetHeader(key, value)`

### Sound

* **Purpose**: Audio playback for sound effects and basic background tracks.
* **Functions**:
* `LoadSound(path)`
* `PlaySound(sound)`
* `SetVolume(sound, level)`
* `StopSound(sound)`

### String

* **Purpose**: Advanced text manipulation, substring extraction, and encoding utilities.
* **Functions**:
* `Base64Encode(str)` / `Base64Decode(str)`
* `Join(list, delimiter)`
* `Left(str, count)`
* `Len(str)`
* `Lower(str)` / `Upper(str)`
* `Mid(str, start, length)`
* `Replace(str, target, replacement)`
* `Right(str, count)`
* `Split(str, delimiter)`

### Struct

* **Purpose**: Structural reflection, key enumeration, and layout serialization.
* **Functions**:
* `FromJson(jsonString, structTypeName)`
* `HasKey(structInstance, keyName)`
* `Keys(structInstance)`
* `ToJson(structInstance)`

### System

* **Purpose**: Process execution, shell spawning, and OS-level control.
* **Functions**:
* `Exec(command, args...)`
* `Exit(code)`
* `GetPID()`
* `Sleep(milliseconds)`
* `Spawn(command, args...)`

### Time

* **Purpose**: Temporal calculations, formatting, and date-time manipulation.
* **Functions**:
* `AddDays(timestamp, days)`
* `Diff(timestampA, timestampB)`
* `Format(timestamp, layout)`
* `Now()`
* `Parse(dateString, layout)`
* `Year(t)`, `Month(t)`, `Day(t)`, `Hour(t)`, `Minute(t)`, `Second(t)`

### XML

* **Purpose**: Structured XML document parsing and querying.
* **Functions**:
* `Parse(xmlString)`
* `Query(node, xpath)`
* `ToString(node)`
