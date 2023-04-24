import { createEffect, createSignal } from 'solid-js'
import NavigationBar from './NavigationBar'

export default function PageNavigator(props) {
  const [navigationBarItems, setNavigationBarItems] = createSignal([])

  // TODO: check if createEffect works when props.activeNavigationBarItem changes
  createEffect(() => {
    if (props.activeNavigationBarItem) {
      setState(props.navigationBarItemPageCombinations, props.activeNavigationBarItem)
    } else {
      setState(
        props.navigationBarItemPageCombinations,
        props.navigationBarItemPageCombinations[0].navigationBarItem
      )
    }
  })

  function setState(navigationBarItemPageCombinations, activeNavigationBarItem) {
    const newNavigationBarItems = new Array(navigationBarItemPageCombinations.length)
    for (let i = 0; i < newNavigationBarItems.length; i++) {
      const navigationBarItem = {
        name: navigationBarItemPageCombinations[i].navigationBarItem,
        active:
          navigationBarItemPageCombinations[i].navigationBarItem === activeNavigationBarItem
            ? true
            : false
      }
      newNavigationBarItems[i] = navigationBarItem
    }
    setNavigationBarItems(newNavigationBarItems)
  }

  function handleChange(e) {
    setState(props.navigationBarItemPageCombinations, e)
  }

  return (
    <div>
      <NavigationBar items={navigationBarItems} onChange={handleChange} />
      <div id="page-wrapper">{props.navigationBarItemPageCombinations[0].page}</div>
    </div>
  )
}
