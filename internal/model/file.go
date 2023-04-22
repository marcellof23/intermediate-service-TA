package model

type File struct {
	Base
	ID           int64
	Filename     string
	OriginalName string
	Client       string
	Size         int64
}

type MemFile struct {
	Name     string // The name of the file.
	RootPath string // The absolute path of the file.
}

type FileDir struct {
	Name        string                 // The name of the current directory we're in.
	RootPath    string                 // The absolute path to this directory.
	Files       map[string]*File       // The list of files in this directory.
	Directories map[string]*Filesystem // The list of directories in this directory.
	Prev        *Filesystem            // a reference pointer to this directory's parent directory.
}

type Filesystem struct {
	*FileDir
}
