package drive

type Drive struct {
	FileSystems map[string]FileSystem
}

type FileSystem interface {
	Put(string, []byte) error
	Get(string) ([]byte, error)
	GetURL(string) (string, error)
	GetSignedURL(string) (string, error)
	Delete(string) error
}

func NewDrive(fileSystems map[string]FileSystem) *Drive {
	return &Drive{
		FileSystems: fileSystems,
	}
}

func (d *Drive) Use(fileSystem string) *FileSystem {
	if fs, ok := d.FileSystems[fileSystem]; ok {
		return &fs
	}

	return nil
}
