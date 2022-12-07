package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var root Filesystem
	var pwd Path

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	var current struct {
		Output strings.Builder
		Execution
	}
	var inputCommand bool

	for _, line := range lines {
		log.Println("line:", line)

		if strings.HasPrefix(line, "$ ") {
			current.Execution.Output = current.Output.String()
			current.Execution = Execution{}
			inputCommand = true
		}

		if inputCommand {
			line = strings.TrimPrefix(line, "$ ")
			parts := strings.Fields(line)
			arg0 := Command(parts[0])
			argv := parts[1:]

			switch arg0 {
			case Command_ls:
				inputCommand = false
				current.Execution.Command.Arg0 = arg0
				current.Execution.Command.Argv = argv
			case Command_cd:
				inputCommand = true
				aocutil.Assertf(len(argv) == 1, "cd: expected 1 argument, got %d", len(argv))
				pwd = pwd.Enter(argv[0])
				log.Printf("cd: pwd = %q", pwd)
			default:
				log.Fatalf("unknown command: %v", arg0)
			}

			continue
		}

		switch current.Execution.Command.Arg0 {
		case Command_ls:
			if strings.HasPrefix(line, "dir ") {
				name := strings.TrimPrefix(line, "dir ")
				log.Printf("mkdir: %q", pwd)
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
		maxSizePaths := walkMaxSize(&root.Folder, 100_000)
		fmt.Println("max size paths:", maxSizePaths)

		var maxSizePathSum int64
		for _, size := range maxSizePaths {
			maxSizePathSum += size
		}

		fmt.Println("max size path sum:", maxSizePathSum)
	}

	{
		op := DeleteOp{
			DiskAvailable: 70_000_000,
			DiskMinUnused: 30_000_000,
		}

		path, ok := op.Do(&root.Folder)
		if !ok {
			fmt.Println("found no path to delete")
			return
		}

		fmt.Println("found path to delete:", path)
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

		sum := file.Size()
		if sum <= threshold {
			pwd := path.Enter(file.Name())
			out[pwd.String()] = sum
		}

		walkMaxSizeRec(folder, threshold, path.Enter(file.Name()), out)
	}
}

type DeleteOp struct {
	DiskAvailable int64
	DiskMinUnused int64

	currentSize  int64
	requiredSize int64
	candidates   map[string]int64
}

func (op DeleteOp) Do(dir *Folder) (deleted Path, ok bool) {
	op.currentSize = dir.Size()
	op.requiredSize = op.DiskAvailable - op.DiskMinUnused
	op.candidates = make(map[string]int64)

	var pwd Path
	op.doRec(dir, pwd)
	spew.Dump(op.candidates)

	if len(op.candidates) == 0 {
		return Path{}, false
	}

	type cand struct {
		path string
		size int64
	}

	candidates := make([]cand, 0, len(op.candidates))
	for path, size := range op.candidates {
		candidates = append(candidates, cand{path, size})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].size < candidates[j].size
	})

	candidate := candidates[0]
	spew.Dump(candidate)
	deleted = deleted.Enter(candidate.path)
	return deleted, true
}

func (op *DeleteOp) doRec(dir *Folder, pwd Path) {
	for _, file := range dir.Files {
		folder, ok := file.(*Folder)
		if !ok {
			continue
		}

		pwd := pwd.Enter(file.Name())

		if op.currentSize-folder.Size() <= op.requiredSize {
			log.Println("found candidate:", pwd)
			op.candidates[pwd.String()] = folder.Size()
		}

		op.doRec(folder, pwd)
	}
}

type Execution struct {
	Command struct {
		Arg0 Command
		Argv []string
	}
	Output string
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
		log.Println("mkdir:", part)
		child := root.Stat(part)
		if child == nil {
			child = &Folder{File: File{name: part}}
			log.Printf("mkdir: created %#v", child)
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

func (f *Folder) Delete(name string) bool {
	for i, file := range f.Files {
		if file.Name() == name {
			f.Files = append(f.Files[:i], f.Files[i+1:]...)
			return true
		}
	}
	return false
}

type File struct {
	name string
	size int64
}

func (f *File) Name() string { return f.name }
func (f *File) Size() int64  { return f.size }

type Path []string

func (p Path) Copy() Path { return append(Path(nil), p...) }

func (p Path) Enter(in string) Path {
	if in == ".." {
		return p[:len(p)-1]
	}

	if strings.HasPrefix(in, "/") {
		in = strings.TrimPrefix(in, "/")
		return Path(strings.Split(in, "/")[1:])
	}

	parts := strings.Split(in, "/")
	return append(p.Copy(), parts...)
}

func (p Path) String() string {
	return "/" + strings.Join(p, "/")
}
