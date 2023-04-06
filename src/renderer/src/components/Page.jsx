export default function Page(props) {
  return (
    <div id="max-page-size">
      <h1>{props.title}</h1>
      {props.children}
    </div>
  )
}
