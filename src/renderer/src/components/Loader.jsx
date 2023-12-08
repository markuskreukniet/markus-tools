export default function Loader(props) {
  return <div id="loader" classList={{ displayNone: !props.loading }} />
}
