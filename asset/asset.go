package asset

import (
	"time"

	"github.com/jessevdk/go-assets"
)

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"empty"}, "empty": []string{}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1526216369, 1526216369476487549),
		Data:     nil,
	}, "empty": &assets.File{
		Path:     "empty",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1526216369, 1526216369476487549),
		Data:     nil,
	}}, "")
