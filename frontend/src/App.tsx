import './global.css'
import { LayoutProvider } from './components/layout-provider'
import { RecoilRoot } from 'recoil'
import Main from './components/Main'

export default function App () {
 
  return (
    <RecoilRoot>
      <LayoutProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <Main />
      </LayoutProvider>
    </RecoilRoot>
  )
}


