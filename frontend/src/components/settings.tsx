import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import * as z from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from './ui/form'
import { useEffect, useState } from 'react'
import { useRecoilState } from 'recoil'
import { loadingState } from '@/store/state'
import { GetConfigs, UpdateConfigs, OpenFileDialog} from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'
import { useToast } from '@/components/ui/use-toast'

type Props = {}

const Settings = (props: Props) => {
  
  const [configs, setConfigs] = useState<main.Config[]>([])

  const formSchema = z.object({
    configs: z.object({
      key: z.string(),
      value: z.string().optional()   
    }).array()
  })
  
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  })

  const [, setLoading] = useRecoilState(loadingState)

  const { toast } = useToast()

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    try {
      setLoading(true)
      await UpdateConfigs(data.configs)
      toast({
        description: 'Configs have been updated.',
      })
    } catch(error: any){
      toast({
        description: error,
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    GetConfigs().then(_configs => {
      setConfigs(_configs)
      form.reset({
        configs: _configs
      })
    }).catch((error:any) => {
      toast({
        description: error,
      })
    })
  }, [])

  
  const openFileDialog = async (event: React.MouseEvent<HTMLElement>, i: number) => {
    try {
      const file = await OpenFileDialog()
      if(file) {
        form.setValue(`configs.${i}.value`, file.replaceAll('\\', '/'))
      }
    } catch (error: any) {
      toast({
        description: error
      })
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="max-w-7xl w-full">
        <Card>
          <CardHeader>
            <CardTitle className="text-center">Settings</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid w-full items-start gap-x-10 gap-y-4 sm:grid-cols-2">
              {configs.map((c, i) => {
                return (
                  <div className="flex flex-col space-y-1.5" key={c.key}>
                    <FormField
                      control={form.control}
                      name={`configs.${i}.value` as const}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel htmlFor={c.key}>{c.key}</FormLabel>
                          <FormControl>
                            <Input onClick={(e) => c.key === 'SSH_KEY' ? openFileDialog(e, i): undefined} type={c.key === 'SSH_PASSPHRASE' ? 'password' : 'text'} id={c.key} aria-describedby={c.key} defaultValue={field.value ?? ''} onChange={field.onChange}/>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                )
              })}
            </div>
          </CardContent>
          <CardFooter className="flex justify-end">
            <Button type='submit'>Submit</Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  )
}

export default Settings