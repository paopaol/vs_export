
# vs_export
read visual studio 15/17/19/22 sln file,export clang compile_commands.json

```cmd
Usage: vs_export -s <path> -c <configuration>

Where:
            -s   path                        sln filename
            -c   configuration               project configuration,eg Debug|Win32.
                                             default Debug|Win32
```

## example

```cmd
vs_export.exe  -s NYWinHotspot.sln  -c "Debug|x64"
```

this can export a compile_commands.json. the compile_commands.json can used by clangd or ccls or some other cpp language server.
