package main

/*

+------------------------------------------------+
| FilesByHashGroups                              |
+------------------------------------------------+
| + FilesByHash: []FilesByHash                   |
+------------------------------------------------+
| + DidAppendByHash(FileSystemFileExtra): bool   | // TODO: add dependency to FileSystemFileExtra
+------------------------------------------------+
                        ^
                        | 0..*
                        |
+------------------------------------------------+
| «struct» FilesByHash                           |
+------------------------------------------------+
| + Hash: string                                 |
| + FileSystemFilesExtra: []FileSystemFilesExtra |
+------------------------------------------------+
+------------------------------------------------+
                        ^
                        | 0..*
                        |
+------------------------------------------------+
| «struct» FileSystemFileExtra                   |
+------------------------------------------------+
| + Hash: string                                 |
| + FileSystemFile: FileSystemFile               |
+------------------------------------------------+
+------------------------------------------------+
                        ^
                        | 1
                        |
+------------------------------------------------+
| «struct» FileSystemFile                        |
+------------------------------------------------+
| + Data:         string                         |
| + Path:         string                         |
| + FileMetadata: FileMetadata                   |
+------------------------------------------------+
+------------------------------------------------+
                        ^
                        | 1
                        |
+------------------------------------------------+
| «struct» FileMetadata                          |
+------------------------------------------------+
| + Name:         string                         |
| + TimeModified: time.Time                      |
| + FileMetadata: FileMetadata                   |
| + Size:         int64                          |
| + IsDirectory:  bool                           |
+------------------------------------------------+
+------------------------------------------------+

+------------------------------+
| «struct» FileData            |
+------------------------------+
| + Data:         string       |
| + FileMetadata: FileMetadata |
+------------------------------+
+------------------------------+
               ^
               | 1
               |
+-------------------------------+
| «struct» FileMetadata         |
+-------------------------------+
| + Name:          string       |
| + Path:          string       |
| + DirectoryPath: string       |
| + Hash:          string       |
| + TimeModified:  time.Time    |
| + FileMetadata:  FileMetadata |
| + Size:          int64        |
| + IsDirectory:   bool         |
+-------------------------------+
+-------------------------------+

*/
