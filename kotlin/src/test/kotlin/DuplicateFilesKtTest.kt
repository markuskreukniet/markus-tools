import org.example.getDuplicateFilesAsNewlineSeparatedString
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import utils.deleteDirectoryTrees
import utils.writeFilesByMultipleInputs
import java.io.File
import java.nio.file.Path
import java.nio.file.Paths

class DuplicateFilesKtTest {
  private lateinit var temporaryDirectories: MutableList<Path>

  @Test
  fun `given inputs when there are duplicate files then return duplicates in newline-separated string`() {
    // arrange
    val contents = arrayOf(
      "content 1\ncontent 1",
      "content 2\ncontent 2",
      "content 3 1\ncontent 3 1"
    )
    val input = """
      empty,,,;
      directory 2/empty,,,;
      directory 1,,txt 1.txt,;
      directory 1,,txt 1 2.txt,${contents[0]};
      directory 2/directory 3,,txt 2-3.txt,${contents[0]};
      directory 2/directory 3,,txt 2-3 2.txt,${contents[1]};
      directory 2/directory 3,,txt 2-3 3.txt,${contents[1]};
      directory 2/directory 4,,txt 2-4.txt,${contents[1]};
      directory 5/directory 6/directory 7,,txt 5-6-7.txt,${contents[2]};
      directory 8,,txt 8.txt,${contents[2]};
    """
    var wantedOutcome = """
      directory 1/txt 1 2.txt
      directory 2/directory 3/txt 2-3.txt

      directory 2/directory 3/txt 2-3 2.txt
      directory 2/directory 3/txt 2-3 3.txt
      directory 2/directory 4/txt 2-4.txt

      directory 5/directory 6/directory 7/txt 5-6-7.txt
      directory 8/txt 8.txt
    """

    val pair = writeFilesByMultipleInputs(input).getOrThrow()
    temporaryDirectories = pair.first ?: fail()
    val inputPaths = pair.second ?: fail()
    val inputPathsArray = inputPaths.toTypedArray()

    val lines = wantedOutcome.lines()
    val trimmedLines = mutableListOf<String>()
    lines.forEach { line ->
      val trimmed = line.trim()
      if (trimmed.isNotEmpty()) {
        // 'Paths.get().toString()'
        // converts a slash-separated path to the native path format for the current operating system.
        trimmedLines.add(Paths.get(trimmed).toString())
      }
    }
    wantedOutcome = trimmedLines.joinToString("\n")

    // act
    var outcome = getDuplicateFilesAsNewlineSeparatedString(inputPathsArray).getOrThrow() ?: fail()

    // assert
    temporaryDirectories.forEach { directory ->
      outcome = outcome.replace("$directory${File.separator}", "")
    }

    outcome.split("\n\n").forEach { substring ->
      wantedOutcome = wantedOutcome.replaceFirst(substring, "")
    }
    wantedOutcome = wantedOutcome.replace("\n\n", "")

    assertEquals(wantedOutcome, "")
  }

  @AfterEach
  fun tearDown() {
    deleteDirectoryTrees(temporaryDirectories)
  }
}
