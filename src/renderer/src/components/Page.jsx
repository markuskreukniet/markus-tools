export default function Page(props) {
  return (
    <div>
      <h1>{props.title}</h1>
      {props.children}
    </div>
  )
}
