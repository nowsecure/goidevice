package afc

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/afc.h>
// int afc_length(char **arr)
// {
// 		int length = 0;
// 		int k = 0;
// 		for (k = 0; arr[k] != NULL; k++) {
// 			length = length + 1;
// 		}
// 		return length;
// }
import "C"
import (
	"unsafe"

	"github.com/nowsecure/goidevice/idevice"
	"github.com/nowsecure/goidevice/lockdown"
)

type AFC struct {
	a C.afc_client_t
}

func NewClient(device idevice.Device, svc *lockdown.Service) (*AFC, error) {
	var a C.afc_client_t
	err := resultToError(
		C.afc_client_new(
			(C.idevice_t)(idevice.GetPointer(device)),
			(C.lockdownd_service_descriptor_t)(svc.GetDescriptor()),
			&a,
		),
	)
	return &AFC{a}, err
}

func (a *AFC) WalkDirectory(path string) ([]SourceFile, error) {
	dir, err := a.ReadDirectory(path)
	if err != nil {
		return []SourceFile{}, err
	}

	files := []SourceFile{}
	for _, f := range dir {
		if f.format == "S_IFDIR" {
			childDir, err := a.WalkDirectory(f.Name)
			if err != nil {
				return []SourceFile{}, err
			}
			files = append(files, childDir...)
			continue
		}
		files = append(files, f)
	}
	return files, nil
}

type SourceFile struct {
	Name string

	format string
}

func (a *AFC) ReadDirectory(path string) ([]SourceFile, error) {
	var directoryC **C.char
	defer C.afc_dictionary_free(directoryC)

	sourceFiles := []SourceFile{}
	pathC := C.CString(path)
	defer C.free(unsafe.Pointer(pathC))

	err := resultToError(C.afc_read_directory(a.a, pathC, &directoryC))
	if err != nil {
		return []SourceFile{}, err
	}

	directory := unsafe.Slice(directoryC, C.afc_length(directoryC))
	for i := range directory {
		file := SourceFile{
			Name: C.GoString(directory[i]),
		}
		if file.Name == ".." || file.Name == "." {
			continue
		}

		var fileInfoC **C.char
		C.afc_get_file_info(a.a, directory[i], &fileInfoC)
		if fileInfoC != nil {
			fileInfo := unsafe.Slice(fileInfoC, C.afc_length(fileInfoC))

			for j := 0; j < len(fileInfo); j += 2 {
				if C.GoString(fileInfo[j]) == "st_ifmt" {
					file.format = C.GoString(fileInfo[j+1])
				}
			}
		}
		C.afc_dictionary_free(fileInfoC)
		sourceFiles = append(sourceFiles, file)
	}
	return sourceFiles, nil
}
