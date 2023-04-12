package utils

import java.nio.file.Files
import java.nio.file.Path
import kotlin.io.path.deleteIfExists
import kotlin.io.path.exists

fun deleteDirectoryTrees(directoryPaths: MutableList<Path>): Result<Unit> = runCatching {
  directoryPaths.forEach { directoryPath ->
    if (!directoryPath.exists()) {
      return@forEach
    }

    Files.walk(directoryPath)
      .sorted(Comparator.reverseOrder())  // Delete files before directories
      .forEach { path -> path.deleteIfExists() }
  }
}
