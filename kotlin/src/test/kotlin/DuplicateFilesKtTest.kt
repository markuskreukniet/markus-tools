import org.example.getDuplicateFilesAsNewlineSeparatedString
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class DuplicateFilesKtTest {
  @Test
  fun `given nothing when always then result is test`() {
    val result = getDuplicateFilesAsNewlineSeparatedString()
    assertEquals("test", result)
  }
}
