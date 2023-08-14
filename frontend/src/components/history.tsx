
import { ColumnDef } from '@tanstack/react-table'
import { useEffect, useState } from 'react'
import { GetMessages } from '../../wailsjs/go/main/App'
import { DataTable } from './data-table'
import { main } from '../../wailsjs/go/models'
import { Button } from './ui/button'
import { SendReversalMessage } from '../../wailsjs/go/main/App'
import AtmResponseDialog from './atm-response-dialog'
import { useRecoilState } from 'recoil'
import { loadingState, messageState, pageState } from '@/store/state'
import { useToast } from './ui/use-toast'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from '@/components/ui/card'


export function History () {
  const columns: ColumnDef<main.Message>[] = [
    {
      accessorKey: 'id',
      header: () => <div className="text-center">ID</div>,
    },
    {
      accessorKey: 'traceNumber',
      header: () => <div className="text-center">Trace Number</div>,
    },
    {
      accessorKey: 'rrn',
      header: () => <div className="text-center">RRN</div>,
    },
    {
      accessorKey: 'transactionAmount',
      header: () => <div className="text-center">Amount</div>,
    },
    {
      accessorKey: 'transaction',
      header: () => <div className="text-center">Transaction</div>,
    },
    {
      accessorKey: 'action',
      header: () => <div className="text-center">Action</div>,
      cell: ({row}) => {
        const disabled = String(row.getValue('transaction')).includes('REVERSAL')
        return (
          <div className='flex justify-center space-x-5'>
            <Button variant="default" disabled={disabled} onClick={() => loadMessage(row.original)}>Load</Button> 
            <Button onClick={() => sendMessage(row.getValue('id'))} variant="destructive" disabled={disabled}>Revert</Button>
          </div>
        )
      }
    },
  ]

  const [, setLoading] = useRecoilState(loadingState)
  const [messages, setMessages] = useState<main.Message[]>([])
  const [isOpen, setOpen] = useState(false)
  const [atmResponse, setAtmResponse] = useState<main.AtmResponse>()
  const [, setMessage] = useRecoilState(messageState)
  const [, setPage] = useRecoilState(pageState)

  const sendMessage = async (id: number) =>{
    let response
    setLoading(true)
    try {
      response = await SendReversalMessage(id)
      setMessages(await GetMessages(1))
      setAtmResponse(response)
      setOpen(true)
    } catch(error:any) {
      toast({
        description: error,
      })
    } finally  {
      setLoading(false)
    }
    
  }

  const loadMessage = (message: main.Message) => {
    setMessage(message)
    setPage('home')
  }

  const { toast } = useToast()

  useEffect(() => {
    GetMessages(1).then(_messages => {setMessages(_messages)}).catch((error:any) => {
      toast({
        description: error,
      })
    })
  }, [])
   

  return (
    <>
      <AtmResponseDialog isOpen={isOpen} setOpen={setOpen} atmResponse={atmResponse}/>
      <Card className='w-full'>
        <CardHeader>
          <CardTitle className="text-center">History</CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable columns={columns} data={messages} />
        </CardContent>
      </Card>
    </>
  
  )
}