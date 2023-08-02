export default function FilePathSelector(props) {
  async function clickInput() {
    props.onChange(await window.dialog.openFileDialogBE(props.directory))
  }

  return <button onClick={clickInput}>{`add a ${props.directory ? 'directory' : 'file'}`}</button>
}
