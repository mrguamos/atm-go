import { atom } from 'recoil'
import { main } from '../../wailsjs/go/models'

export const loadingState = atom({
  key: 'loading', 
  default: false, 
})

export const pageState = atom({
  key:'page',
  default:'home'
})

export const messageState = atom({
  key: 'message',
  default: new main.Message()
})

export const tunnelState = atom({
  key: 'tunnel',
  default: false
})