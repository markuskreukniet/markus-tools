export default function Loader(props) {
  return <div id="loader" classList={{ 'display-none': !props.loading }} />
}

// TODO: spinner is better naming since it does not load anything