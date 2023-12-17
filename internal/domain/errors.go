package domain

import (
	"fmt"
)

var (
	ErrInvalidLinkFormat         = fmt.Errorf("invalid link format")
	ErrAlbumDownloadsUnsupported = fmt.Errorf("album downloads are not yet supported")
)
