import { createEffect, createSignal } from 'solid-js'
import NavigationBar from './NavigationBar'

export default function PageNavigator(props) {
  const [navigationBarItems, setNavigationBarItems] = createSignal([])

  createEffect(() => {
    const activeNavigationBarItem = props.activeNavigationBarItem
      ? props.activeNavigationBarItem
      : props.items[0]

    const newNavigationBarItems = new Array(props.navigationBarItemPageCombinations.length)
    for (let i = 0; i < newNavigationBarItems.length; i++) {
      const navigationBarItem = {
        name: props.navigationBarItemPageCombinations[i].navigationBarItem,
        active:
          props.navigationBarItemPageCombinations[i].navigationBarItem === activeNavigationBarItem
            ? true
            : false
      }
      newNavigationBarItems[i] = navigationBarItem
    }
    setNavigationBarItems(newNavigationBarItems)
  })

  return (
    <div>
      <NavigationBar items={navigationBarItems} />
      <div id="page-wrapper">{props.navigationBarItemPageCombinations[0].page}</div>
    </div>
  )
}
