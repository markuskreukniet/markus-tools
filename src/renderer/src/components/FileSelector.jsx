export default function FileSelector(props) {
  return (
    <input
      type="file"
      webkitdirectory={props.folder}
      onClick={(e) => (e.target.value = '') /* makes selecting the same file or folder possible */}
      onChange={(e) => props.onChange(e.target.files)}
    />
  )
}
