package paths

import "path/filepath"

const DirName = ".devenv"
const SnapshotFileName = "snapshot.yaml"
const RestorePS1Name = "restore.ps1"
const RestoreSHName = "restore.sh"

func Dir(root string) string {
	return filepath.Join(root, DirName)
}

func SnapshotPath(root string) string {
	return filepath.Join(Dir(root), SnapshotFileName)
}

func RestorePS1Path(root string) string {
	return filepath.Join(Dir(root), RestorePS1Name)
}

func RestoreSHPath(root string) string {
	return filepath.Join(Dir(root), RestoreSHName)
}

