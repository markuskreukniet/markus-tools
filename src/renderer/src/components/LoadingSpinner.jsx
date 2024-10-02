export default function LoadingSpinner(props) {
  return <div id="loading-spinner" classList={{ 'display-none': !props.loading }} />
}
