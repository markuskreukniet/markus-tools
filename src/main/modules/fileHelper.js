export default function isNotAZeroByteFile(stats) {
  return stats.size > 0
}
