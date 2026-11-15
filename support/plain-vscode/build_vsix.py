import json
import os
import tempfile
import zipfile
from pathlib import Path
from xml.sax.saxutils import escape, quoteattr


EXTENSION_DIR = Path(__file__).resolve().parent
PACKAGE_PATH = EXTENSION_DIR / "package.json"


def build_manifest(package):
    categories = ", ".join(package.get("categories", []))
    engine = package["engines"]["vscode"]
    return f"""<?xml version="1.0" encoding="utf-8"?>
<PackageManifest Version="2.0.0" xmlns="http://schemas.microsoft.com/developer/vsx-schema/2011">
  <Metadata>
    <Identity Language="en-US" Id={quoteattr(package["name"])} Version={quoteattr(package["version"])} Publisher={quoteattr(package["publisher"])} />
    <DisplayName>{escape(package["displayName"])}</DisplayName>
    <Description xml:lang="en-US">{escape(package["description"])}</Description>
    <Categories>{escape(categories)}</Categories>
    <Properties>
      <Property Id="Microsoft.VisualStudio.Code.Engine" Value={quoteattr(engine)} />
    </Properties>
  </Metadata>
  <Installation>
    <InstallationTarget Id="Microsoft.VisualStudio.Code" />
  </Installation>
  <Dependencies />
  <Assets>
    <Asset Type="Microsoft.VisualStudio.Code.Manifest" Path="extension/package.json" />
  </Assets>
</PackageManifest>
"""


def extension_files(package):
    files = [PACKAGE_PATH]
    grammar_paths = package.get("contributes", {}).get("grammars", [])
    for grammar in grammar_paths:
        relative_path = Path(grammar["path"])
        if relative_path.is_absolute() or ".." in relative_path.parts:
            raise ValueError(f"Invalid extension asset path: {grammar['path']}")
        files.append(EXTENSION_DIR / relative_path)
    return files


def main():
    with PACKAGE_PATH.open(encoding="utf-8") as package_file:
        package = json.load(package_file)

    output_path = EXTENSION_DIR / f'{package["name"]}-{package["version"]}.vsix'
    with tempfile.TemporaryDirectory(dir=EXTENSION_DIR) as temporary_dir:
        temporary_path = Path(temporary_dir) / output_path.name
        with zipfile.ZipFile(temporary_path, "w", zipfile.ZIP_DEFLATED) as archive:
            archive.writestr(
                "[Content_Types].xml",
                """<?xml version="1.0" encoding="utf-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="json" ContentType="application/json" />
  <Default Extension="xml" ContentType="text/xml" />
  <Default Extension="vsixmanifest" ContentType="text/xml" />
</Types>
""",
            )
            archive.writestr("extension.vsixmanifest", build_manifest(package))
            for source_path in extension_files(package):
                if not source_path.is_file():
                    raise FileNotFoundError(f"Extension asset not found: {source_path}")
                archive.write(
                    source_path,
                    f"extension/{source_path.relative_to(EXTENSION_DIR).as_posix()}",
                )
        os.replace(temporary_path, output_path)

    print(f"Created {output_path}")


if __name__ == "__main__":
    main()
