package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var root Filesystem
	var pwd Path
	var cmd Command

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	for _, line := range lines {
		if strings.HasPrefix(line, "$ ") {
			line = strings.TrimPrefix(line, "$ ")
			parts := strings.Fields(line)

			cmd = Command(parts[0])
			argv := parts[1:]

			switch cmd {
			case Command_cd:
				aocutil.Assertf(len(argv) == 1, "cd: expected 1 argument, got %d", len(argv))
				pwd = pwd.Enter(argv[0])
			}

			continue
		}

		switch cmd {
		case Command_ls:
			if strings.HasPrefix(line, "dir ") {
				name := strings.TrimPrefix(line, "dir ")
				root.Mkdir(pwd.Enter(name))
				continue
			}

			// Probably a file. Parse its size.
			var size int64
			var name string
			_, err := fmt.Sscanf(line, "%d %s", &size, &name)
			aocutil.E1(err)

			dir := root.Mkdir(pwd)
			dir.Touch(name, size)
		}
	}

	{
		fmt.Print("Part 1: ")

		maxSizePaths := walkMaxSize(&root.Folder, 100_000)

		var maxSizePathSum int64
		for _, size := range maxSizePaths {
			maxSizePathSum += size
		}

		fmt.Println("max size path sum:", maxSizePathSum)
	}

	{
		fmt.Print("Part 2: ")

		op := DeleteOp{
			DiskAvailable: 70_000_000,
			DiskMinUnused: 30_000_000,
		}

		path, size, ok := op.Do(&root.Folder)
		if !ok {
			fmt.Println("found no path to delete")
			return
		}

		fmt.Printf("found path to delete: %s (%d)", path, size)
		fmt.Println()
	}
}

func walkMaxSize(f *Folder, threshold int64) map[string]int64 {
	out := make(map[string]int64)
	walkMaxSizeRec(f, threshold, Path{}, out)
	return out
}

func walkMaxSizeRec(f *Folder, threshold int64, path Path, out map[string]int64) {
	for _, file := range f.Files {
		folder, ok := file.(*Folder)
		if !ok {
			continue
		}

		pwd := path.Enter(file.Name())
		sum := file.Size()
		if sum <= threshold {
			out[pwd.String()] = sum
		}

		walkMaxSizeRec(folder, threshold, pwd, out)
	}
}

type DeleteOp struct {
	DiskAvailable int64
	DiskMinUnused int64

	currentSize  int64
	requiredSize int64
	candidates   map[string]int64
}

func (op DeleteOp) Do(dir *Folder) (Path, int64, bool) {
	op.currentSize = dir.Size()
	op.requiredSize = op.DiskAvailable - op.DiskMinUnused
	op.candidates = make(map[string]int64)

	var pwd Path
	op.doRec(dir, pwd)

	if len(op.candidates) == 0 {
		return Path{}, 0, false
	}

	candidates := aocutil.MapPairs(op.candidates)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].V < candidates[j].V
	})

	candidate := candidates[0]
	return ParsePath(candidate.K), candidate.V, true
}

func (op *DeleteOp) doRec(dir *Folder, pwd Path) {
	for _, file := range dir.Files {
		folder, ok := file.(*Folder)
		if !ok {
			continue
		}

		pwd := pwd.Enter(file.Name())

		if op.currentSize-folder.Size() <= op.requiredSize {
			op.candidates[pwd.String()] = folder.Size()
		}

		op.doRec(folder, pwd)
	}
}

type Command string

const (
	Command_ls Command = "ls"
	Command_cd Command = "cd"
)

type Filesystem struct {
	Folder
}

func (fs *Filesystem) Mkdir(fsPath Path) *Folder {
	root := &fs.Folder
	var ok bool

	for _, part := range fsPath {
		child := root.Stat(part)
		if child == nil {
			child = &Folder{File: File{name: part}}
			root.Files = append(root.Files, child)
		}
		root, ok = child.(*Folder)
		aocutil.Assertf(ok, "%q: expected folder, got %T", fsPath, child)
	}

	return root
}

type Filer interface {
	Name() string
	Size() int64
}

type Folder struct {
	File
	Files []Filer
}

func (f *Folder) Stat(name string) Filer {
	for _, file := range f.Files {
		if file.Name() == name {
			return file
		}
	}
	return nil
}

func (f *Folder) Touch(name string, size int64) {
	if file := f.Stat(name); file != nil {
		return
	}

	file := &File{name, size}
	f.Files = append(f.Files, file)
}

func (f *Folder) Size() int64 {
	var size int64
	for _, file := range f.Files {
		size += file.Size()
	}
	return size
}

type File struct {
	name string
	size int64
}

func (f *File) Name() string { return f.name }
func (f *File) Size() int64  { return f.size }

type Path []string

func ParsePath(str string) Path {
	var path Path
	return path.Enter(str)
}

func (p Path) Copy() Path { return append(Path(nil), p...) }

func (p Path) Enter(in string) Path {
	if in == ".." {
		return p[:len(p)-1]
	}

	if strings.HasPrefix(in, "/") {
		if in == "/" {
			return Path{}
		}
		in = strings.TrimPrefix(in, "/")
		return Path(strings.Split(in, "/"))
	}

	parts := strings.Split(in, "/")
	return append(p.Copy(), parts...)
}

func (p Path) String() string {
	return "/" + strings.Join(p, "/")
}
