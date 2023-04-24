import { For } from 'solid-js'

export default function NavigationBar(props) {
  return (
    <ul id="navigation-bar">
      <For each={props.items()}>
        {(item) => (
          <li
            onClick={() => props.onChange(item.name)}
            classList={{ navigationBarItemActive: item.active }}
          >
            {item.name}
          </li>
        )}
      </For>
    </ul>
  )
}
