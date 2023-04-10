export default function Page(props) {
  return (
    <div id="page">
      <h1>{props.title}</h1>
      {props.children}
    </div>
  )
}
