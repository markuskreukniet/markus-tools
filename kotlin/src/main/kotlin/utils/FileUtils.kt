package org.example.utils

import java.io.File
import java.net.URLConnection
import java.nio.file.Path
import java.nio.file.attribute.FileTime
import java.time.Instant
import kotlin.io.path.getLastModifiedTime

// The F prefix of a class name means feature.

interface DuplicateFileInfo {
  val file: File
  val size: Long
}

data class FDuplicateFilesFileInfo(
  override val file: File,
  override val size: Long,
) : DuplicateFileInfo

data class FDateRangeFileInfo(
  override val file: File,
  override val size: Long,
  val path: Path,
  val timeModified: Instant,
  var newName: String? // We need the 'newName' property because we cannot change the name of a File instance directly.
) : DuplicateFileInfo

data class FTextFilesFileInfo(
  val file: File,
  val absolutePath: Path
)

data class CompleteFileInfo(
  val file: File,
  val size: Long,
  var absolutePath: Path,
  var absoluteDirectoryPath: Path,
  val timeModified: FileTime?,
)

data class FileData(
  val content: String,
  val completeFileInfo: CompleteFileInfo
)

enum class FileFilterMode {
  FILES,
  NON_ZERO_BYTE_FILES,
  FILES_AND_DIRECTORIES,
  NON_ZERO_BYTE_FILES_AND_DIRECTORIES,
  DIRECTORIES
}

enum class FileType {
  ALL_FILES,
  TEXT_FILES
}

fun isTextFile(file: File): Result<Boolean> = runCatching {
  val mimeType = file.inputStream().use { inputStream ->
    URLConnection.guessContentTypeFromStream(inputStream)
  }

  mimeType?.startsWith("text") == true
}

fun resolveDirectoryPath(filePath: Path, isDirectory: Boolean): Path {
  return if (isDirectory || filePath.parent == null) {
    filePath
  } else {
    filePath.parent
  }
}

fun filterAndHandleFileInfo(
  file: File, mode: FileFilterMode, type: FileType, absoluteFilePath: Path, handler: (CompleteFileInfo) -> Unit
): Result<Unit> = runCatching {
  val size = if (file.isFile) file.length() else 0L

  // is file check
  if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
    return@runCatching
  }

  // is directory check
  if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
    return@runCatching
  }

  // is zero byte file check
  if (file.isFile && size == 0L &&
    (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
    return@runCatching
  }

  // is text file check
  if (type == FileType.TEXT_FILES) {
    if (!isTextFile(file).getOrThrow()) {
      return@runCatching
    }
  }

  handler(CompleteFileInfo(
    file = file,
    size = file.length(),
    absolutePath = absoluteFilePath,
    absoluteDirectoryPath = resolveDirectoryPath(absoluteFilePath, file.isDirectory),
    timeModified = absoluteFilePath.getLastModifiedTime()
  ))
}

fun walkFilterAndHandleFileInfo(
  absoluteFilePath: Path,
  mode: FileFilterMode,
  type: FileType,
  handler: (CompleteFileInfo) -> Unit
): Result<Unit> = runCatching {
  val rootFile = absoluteFilePath.toFile()

  if (!rootFile.isFile && !rootFile.isDirectory) {
    return@runCatching
  }

  val files = if (rootFile.isDirectory) rootFile.walk() else sequenceOf(rootFile)
  files.forEach { file ->
    filterAndHandleFileInfo(file, mode, type, absoluteFilePath, handler).onFailure { throw it }
  }
}
