import { For } from 'solid-js'

export default function NavigationBar(props) {
  return (
    <ul id="navigation-bar">
      <For each={props.items()}>
        {(item) => (
          <li
            onMouseDown={() => props.onChange(item.name)}
            classList={{ 'navigation-bar__item--active': item.active }}
          >
            {item.name}
          </li>
        )}
      </For>
    </ul>
  )
}
