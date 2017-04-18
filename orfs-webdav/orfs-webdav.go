package main

import (
	"github.com/cetex/ORFS/orfs"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
)

type fs struct {
	orfs *orfs.Orfs
}

func (f *fs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return f.orfs.Mkdir(name, perm)
}

func (f *fs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return f.orfs.OpenFile(name, flag, perm)
}

func (f *fs) RemoveAll(ctx context.Context, name string) error {
	return f.orfs.RemoveAll(name)
}

func (f *fs) Rename(ctx context.Context, oldName, newName string) error {
	return f.orfs.Rename(oldName, newName)
}

func (f *fs) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return f.orfs.Stat(name)
}

func main() {
	c := &fs{orfs.NewORFS("test", "test_metadata")}
	c.orfs.SetDebugLog(os.Stdout)
	c.orfs.SetLog(os.Stderr)
	if err := c.orfs.Connect(); err != nil {
		panic(err)
	}
	srv := &webdav.Handler{
		Prefix:     "/",
		FileSystem: c,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.Printf("WEBDAV: %#s, ERROR: %v", r, err)
		},
	}
	http.Handle("/", srv)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}
}
