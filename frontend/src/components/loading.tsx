'use client'

import { Loader } from 'lucide-react'

export default function Loading() {
  return <div className="flex fixed w-full h-screen justify-center items-center top-0 left-0 z-50 bg-background/80 backdrop-blur-sm"><Loader className="animate-spin"/></div>
}