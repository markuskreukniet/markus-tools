export default function Spacer(props) {
  return <div classList={`${props.width || 'fullWidth'} ${props.height || 'fullHeight'}`} />
}
