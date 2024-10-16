import org.example.getDuplicateFilesAsNewlineSeparatedString
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import java.nio.file.Paths

class DuplicateFilesKtTest {
  @Test
  fun `given nothing when always then result is test`() {
    val paths = arrayOf(
      Paths.get("/path1/path2"),
      Paths.get("/path3")
    )

    val result = getDuplicateFilesAsNewlineSeparatedString(paths).getOrThrow() ?: return
    assertEquals("test", result)
  }
}
