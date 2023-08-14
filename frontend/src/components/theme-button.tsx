'use client'

import { Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useTheme } from '@/components/layout-provider'

export function ThemeButton() {
  const { theme, setTheme } = useTheme()

  return (
    <>
      <Button variant="ghost" size="icon" onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}>
        <Sun className="h-[1.2rem] w-[1.2rem] dark:-rotate-0 dark:scale-100 transition-all rotate-90 scale-0 dark:z-10"/>
        <Moon className="absolute h-[1.2rem] w-[1.2rem] dark:-rotate-90 dark:scale-0 transition-all rotate-0 scale-100 dark:-z-10" />
      </Button>
   
    </>
  )
}
