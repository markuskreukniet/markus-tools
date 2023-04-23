import { createEffect } from 'solid-js'
import NavigationBar from './NavigationBar'

export default function PageNavigator(props) {
  const navigationBarItems = []

  createEffect(() => {
    for (const combination of props.navigationBarItemPageCombinations) {
      navigationBarItems.push(combination.navigationBarItem)
    }
  })

  return (
    <div>
      <NavigationBar items={navigationBarItems} />
      <div id="page-wrapper">{props.navigationBarItemPageCombinations[0].page}</div>
    </div>
  )
}
