package webdav

func Mount(mountUrl, location string) error {
	return execMountCommand("net", "use", "*", mountUrl, location)
}
