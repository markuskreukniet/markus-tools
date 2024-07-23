package main

/*

+-----------------------------+
| FileSystemFile              |
|-----------------------------|
| + Data         string       |
| + Path         string       |
| + FileMetadata FileMetadata |
+-----------------------------+
|                             |
+-----------------------------+
               ^
               | 1
               |
+-----------------------------+
| FileMetadata                |
|-----------------------------|
| + Name         string       |
| + TimeModified time.Time    |
| + FileMetadata FileMetadata |
| + Size         int64        |
| + IsDirectory  bool         |
+-----------------------------+
|                             |
+-----------------------------+

*/
