import { For } from 'solid-js'

export default function NavigationBar(props) {
  return (
    <ul id="navigation-bar">
      <For each={props.items()}>
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
