package main

/*

+------------------------------------------------+
| FilesByHashGroups                              |
+------------------------------------------------+
| + FilesByHashGroup: []FilesByHashGroup         |
+------------------------------------------------+
| + DidAppendByHash(FileSystemFileExtra): bool   | // TODO: add dependency to FileSystemFileExtra
+------------------------------------------------+
                        ^
                        | 0..*
                        |
+------------------------------------------------+
| FilesByHashGroup                               |
+------------------------------------------------+
| + Hash: string                                 |
| + FileSystemFilesExtra: []FileSystemFilesExtra |
+------------------------------------------------+
+------------------------------------------------+
                        ^
                        | 0..*
                        |
+------------------------------------------------+
| FileSystemFileExtra                            |
+------------------------------------------------+
| + Hash: string                                 |
| + FileSystemFile: FileSystemFile               |
+------------------------------------------------+
+------------------------------------------------+
                        ^
                        | 1
                        |
+------------------------------------------------+
| FileSystemFile                                 |
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
| FileMetadata                                   |
+------------------------------------------------+
| + Name:         string                         |
| + TimeModified: time.Time                      |
| + FileMetadata: FileMetadata                   |
| + Size:         int64                          |
| + IsDirectory:  bool                           |
+------------------------------------------------+
+------------------------------------------------+

*/
