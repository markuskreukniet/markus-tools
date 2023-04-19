import { createEffect, createSignal, For } from 'solid-js'

export default function NavigationBar(props) {
  const [items, setItems] = createSignal([])
  createEffect(() => {
    const activeItem = props.activeItem ? props.activeItem : props.items[0]

    const newItems = new Array(props.items.length)
    for (let i = 0; i < newItems.length; i++) {
      const newItem = {
        name: props.items[i],
        active: props.items[i] === activeItem ? true : false
      }
      newItems[i] = newItem
    }
    setItems(newItems)
  })

  return (
    <ul id="navigation-bar">
      <For each={items()}>
        {(item) => (
          <li
            onClick={() => console.log('set active css class')}
            classList={{ navigationBarItemActive: item.active }}
          >
            {item.name}
          </li>
        )}
      </For>
    </ul>
  )
}
