import { createEffect, createSignal } from 'solid-js'

export default function Loader(props) {
  const [loading, setLoading] = createSignal(false)
  createEffect(() => {
    setLoading(props.loading)
  })

  return <div id="loader" classList={{ displayNone: !loading() }} />
}
