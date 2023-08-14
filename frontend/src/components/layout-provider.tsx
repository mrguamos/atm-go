import { createContext, useContext, useEffect, useState } from 'react'

import Loading from '@/components/loading'
import { useRecoilValue } from 'recoil'
import { loadingState } from '@/store/state'

type LayoutProviderProps = {
  children: React.ReactNode
  defaultTheme?: string
  storageKey?: string
}

type LayoutProviderState = {
  theme: string
  setTheme: (theme: string) => void
}

const initialState = {
  theme: 'system',
  setTheme: () => null
}

const LayoutProviderContext = createContext<LayoutProviderState>(initialState)

export function LayoutProvider ({
  children,
  defaultTheme = 'system',
  storageKey = 'vite-ui-theme',
  ...props
}: LayoutProviderProps) {
  const [theme, setTheme] = useState(
    () => localStorage.getItem(storageKey) || defaultTheme
  )

  useEffect(() => {
    const root = window.document.documentElement

    root.classList.remove('light', 'dark')

    if (theme === 'system') {
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)')
        .matches
        ? 'dark'
        : 'light'

      root.classList.add(systemTheme)
      return
    }

    root.classList.add(theme)
  }, [theme])

  const value = {
    theme,
    setTheme: (theme: string) => {
      localStorage.setItem(storageKey, theme)
      setTheme(theme)
    }
  }

  const loading = useRecoilValue(loadingState)

  return (
    <LayoutProviderContext.Provider {...props} value={value}>
      {loading && <Loading /> }
      {children}
    </LayoutProviderContext.Provider>
  )
}

export const useTheme = () => {
  const context = useContext(LayoutProviderContext)

  if (context === undefined) { throw new Error('useTheme must be used within a LayoutProvider') }

  return context
}
