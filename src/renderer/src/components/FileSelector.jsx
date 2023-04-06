export default function FileSelector(props) {
  function clickInput(e) {
    const input = e.srcElement.parentElement.getElementsByTagName('input')[0]
    input.click()
  }

  return (
    <div>
      <input
        type="file"
        webkitdirectory={props.folder}
        onClick={
          (e) => (e.target.value = '') /* makes selecting the same file or folder possible */
        }
        onChange={(e) => props.onChange(e.target.files)}
        class="display-none"
      />
      <button onClick={clickInput}>{`add a ${props.folder ? 'folder' : 'file'}`}</button>
    </div>
  )
}
