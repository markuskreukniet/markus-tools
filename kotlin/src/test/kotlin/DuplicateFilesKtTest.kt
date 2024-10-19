import org.example.getDuplicateFilesAsNewlineSeparatedString
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import utils.writeFilesByMultipleInputs
import java.nio.file.Paths
import kotlin.io.path.exists

class DuplicateFilesKtTest {
  @Test
  fun `given inputs when there are duplicate files then return duplicates in newline-separated string`() {
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
      directory 8,,txt 8.txt,${contents[2]}
    """
    val wantedOutcome = """
      directory 1\txt 1 2.txt
      directory 2\directory 3\txt 2-3.txt

      directory 2\directory 3\txt 2-3 2.txt
      directory 2\directory 3\txt 2-3 3.txt
      directory 2\directory 4\txt 2-4.txt

      directory 5\directory 6\directory 7\txt 5-6-7.txt
      directory 8\txt 8.txt
    """

    val pair = writeFilesByMultipleInputs(input).getOrThrow()

    val temporaryDirectories = pair.first ?: fail()

//    temporaryDirectories.forEach { directory ->
//      if (directory.exists()) {
//        //
//      }
//    }

    val paths = arrayOf(
      Paths.get("/path1/path2"),
      Paths.get("/path3")
    )

    val result = getDuplicateFilesAsNewlineSeparatedString(paths).getOrThrow() ?: return
    assertEquals("test", result)
  }
}
