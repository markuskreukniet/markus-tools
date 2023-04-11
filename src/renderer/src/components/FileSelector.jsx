export default function FileSelector(props) {
  function clickInput(e) {
    const input = e.srcElement.parentElement.getElementsByTagName('input')[0]
    input.click()
  }

  // TODO: Cancel event on input type="file": https://stackoverflow.com/questions/34855400/cancel-event-on-input-type-file
  // We can then show a loader when clicking the input and dismiss the loader on the cancel event.
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
