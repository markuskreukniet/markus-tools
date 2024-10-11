import org.example.getDuplicateFilesAsNewlineSeparatedString
import org.example.utils.FileSystemNode
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class DuplicateFilesKtTest {
  @Test
  fun `given nothing when always then result is test`() {
    val nodes = arrayOf(
      FileSystemNode("/path1/path2", false),
      FileSystemNode("/path3", false)
    )

    val result = getDuplicateFilesAsNewlineSeparatedString(nodes).getOrThrow() ?: return
    assertEquals("test", result)
  }
}
