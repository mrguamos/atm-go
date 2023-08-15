import { useRecoilState } from 'recoil'
import { ThemeButton } from './theme-button'
import { Landmark, Cable } from 'lucide-react'
import { loadingState, pageState, tunnelState } from '@/store/state'
import { Button } from './ui/button'
import { CloseTunnel, PingTunnel, UseTunnel } from '../../wailsjs/go/main/App'
import { useToast } from '@/components/ui/use-toast'
import { EventsOn } from '../../wailsjs/runtime'
import { useEffect } from 'react'


const Navigation = () => {

  const [, setPage] = useRecoilState(pageState)
  const [tunnel, setTunnel] = useRecoilState(tunnelState)
  const [, setLoading] = useRecoilState(loadingState)
  const { toast } = useToast()
  EventsOn('tunnel', (tunnel: boolean, withToast: boolean = true) => {
    setTunnel(tunnel)
    if(withToast) {
      toast({
        description: `Tunnel has been ${tunnel ? 'enabled.' : 'disabled.'}`,
      })
    }
  })

  useEffect(() => {
    PingTunnel().then()
  }, [])
  

  return (
    <div className="flex min-h-[60px] px-10 sticky backdrop-blur-lg top-0 items-center bg-secondary/50 justify-between">
      <div className='flex items-center'>
        <Landmark className='mr-20'/>
        <div className='flex space-x-10'>
          <span className='text-xl hover:cursor-pointer' onClick={()=> setPage('home')}>ATM</span>
          <span className='text-xl hover:cursor-pointer' onClick={()=> setPage('history')}>HISTORY</span>
          <span className='text-xl hover:cursor-pointer' onClick={()=> setPage('settings')}>SETTINGS</span>
        </div>
      </div>
      <div className='flex items-center'>
        <Button variant="ghost" size="icon" title='Tunnel' onClick={async () => {
          try {
            setLoading(true)
            if(!tunnel) {
              await UseTunnel()
            }else {
              await CloseTunnel()
            }
          } catch(error: any) {
            toast({
              description: error,
            })
          } finally {
            setLoading(false)   
          }
        }}>
          <Cable color={tunnel ? 'green': 'red'} />
        </Button>
        <ThemeButton />
      </div>
    </div>
  )
}

export default Navigation