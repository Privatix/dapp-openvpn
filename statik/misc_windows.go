package statik

//go:generate cmd /C IF EXIST "statik.go" (del /F /Q statik.go)
//go:generate statik -f -src=. -dest=..
