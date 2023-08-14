import { keyState, pageState } from '@/store/state'
import { useRecoilValue } from 'recoil'
import { AtmForm } from './atm-form'
import Navigation from './navigation'
import { Toaster } from './ui/toaster'
import { History } from './history'
import Config from './settings'

const Main = () => {

  const page = useRecoilValue(pageState)
  const key = useRecoilValue(keyState)
  return (
    <div className="flex flex-col min-h-screen">
      <Navigation />
      <div className="flex grow container mx-auto px-4 py-20">
        {page === 'home' && <AtmForm key={key}/>}
        {page === 'history' && <History />}
        {page === 'settings' && <Config />}
        <Toaster />
      </div>
    </div>
  )
}

export default Main