# Armadillo (Windows version)
## Task
Develop a program that prevents the current directory (where it resides) from creating, copying, or renaming files with specified names (you can use file masks). 
Store the list of names or their templates in the template.tbl file as text. 
This file shall be protected from deletion, unauthorized viewing and modification. 
If you install the program, you can use the password stored in the first line of the template.tbl file to disable it. 
The program must turn on and off the protection mode.

## Build
```bash
go build -o armadillo.exe ./cmd/main.go
```

## Run
List of commands:
* install - add armadillo to Windows Service
* remove - remove armadillo from Windows Service
* start - boot armadillo and turn on protection (requires password)
* stop - stopping armadillo, if password is correct

## Sources
[Repository with code, which Manupulating with ACL](https://github.com/hectane/go-acl)
