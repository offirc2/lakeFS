package webdav

func Mount(mountUrl, location string) error {
	return execMountCommand("mount_webdav", "-S", mountUrl, location)
}

func Umount(location string) error {
	return execMountCommand("umount", location)
}
