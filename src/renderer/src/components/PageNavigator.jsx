import { createEffect, createSignal } from 'solid-js'
import NavigationBar from './NavigationBar'

export default function PageNavigator(props) {
  const [navigationBarItems, setNavigationBarItems] = createSignal([])
  const [activePage, setActivePage] = createSignal(<></>)

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
        active: false
      }
      if (navigationBarItemPageCombinations[i].navigationBarItem === activeNavigationBarItem) {
        navigationBarItem.active = true

        setActivePage(navigationBarItemPageCombinations[i].page)
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
      <div id="page-navigator__page-wrapper">{activePage()}</div>
    </div>
  )
}
