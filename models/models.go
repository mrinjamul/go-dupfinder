package models

type FileInfo struct {
	FileName  string
	FileSize  int64
	ChunkHash string
	Hash      string
}
