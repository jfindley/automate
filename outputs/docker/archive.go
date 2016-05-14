package docker

import (
	"archive/tar"
	"bytes"
	"os"
    "io"

	"github.com/jfindley/testfs"
)

func archive(fs testfs.FileSystem) (out []byte, err error) {
    buf := new(bytes.Buffer)
	t := tar.NewWriter(buf)
    
    origDir, err := fs.Getwd()
    if err != nil {
		return
	}

	// Start at the root of the fs
	fs.Chdir("/")

	err = walk(fs, t)
	if err != nil {
		return
	}
    
    err = fs.Chdir(origDir)
    if err != nil {
		return
	}

	err = t.Close()
    
	return buf.Bytes(), err
}

func walk(fs testfs.FileSystem, t *tar.Writer) (err error) {

	cwd, err := fs.Getwd()
	if err != nil {
		return
	}

	f, err := fs.OpenFile(cwd, os.O_RDONLY, 0)
	if err != nil {
		return
	}

	fi, err := f.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range fi {
		// Add a new file to the archive
		var hdr *tar.Header
		hdr, err = tarHdr(file)
		if err != nil {
			return
		}
        
        // Make header name fully qualified
		if cwd[len(cwd)-1] == '/' {
			hdr.Name = cwd + hdr.Name
		} else {
			hdr.Name = cwd + "/" + hdr.Name
		}
        
        // Strip leading slashes
        if hdr.Name[0] == '/' {
            hdr.Name = hdr.Name[1:]
        }
        
		err = t.WriteHeader(hdr)
		if err != nil {
			return
		}
        
		// Write the file data to the archive
		if !file.IsDir() && file.Size() > 0 {
			var target testfs.File

			target, err = fs.Open(file.Name())
			if err != nil {
				return
			}

			_, err = io.Copy(t, target)
            if err != nil {
				return
			}
		}

		// Close the file in the archive
		err = t.Flush()
		if err != nil {
			return
		}

		// If a directory, recurse into it
		if file.IsDir() {
			err = fs.Chdir(file.Name())
			if err != nil {
				return
			}

			return walk(fs, t)
		}

	}

	return
}

func tarHdr(file os.FileInfo) (hdr *tar.Header, err error) {

	// We only deal with testfs filesystems.
	sys := file.Sys().(*testfs.Stat_t)

	hdr, err = tar.FileInfoHeader(file, sys.Linkname)
	if err != nil {
		return
	}

	// archive.tar doesn't understand testfs, so set this manually.
	hdr.Uid = int(sys.Uid)
	hdr.Gid = int(sys.Gid)
	hdr.Xattrs = sys.Xattrs

	if sys.Linkname != "" {
		hdr.Size = 0
		hdr.Linkname = sys.Linkname
		hdr.Typeflag = '1'
	}

	return
}
