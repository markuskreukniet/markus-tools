import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import utils.deleteDirectoryTrees
import java.nio.file.Path

class FilesToDateRangeDirectoryKtTest {
  private lateinit var temporaryDirectories: MutableList<Path>

  @Test
  fun `given - when - then -`() {
    // arrange
    // var destinationInput = ""

    // val destination = writeFilesBySingleInput(destinationInput).getOrThrow()

    // act
    // val outcome = filesToDateRangeDirectory(inputPathsArray, ).getOrThrow() ?: fail()

    // assert
  }

  @AfterEach
  fun tearDown() {
    deleteDirectoryTrees(temporaryDirectories)
  }
}
