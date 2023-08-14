import { AlertDialog, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from './ui/alert-dialog'
import { main } from '../../wailsjs/go/models'

type Props = {
    atmResponse?: main.AtmResponse
    isOpen: boolean
    setOpen: (isOpen: boolean) => void
}

const AtmResponseDialog = ({atmResponse, isOpen, setOpen}: Props) => {

  return (
    <AlertDialog open={isOpen} onOpenChange={() => setOpen(!isOpen)}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>ATM Response</AlertDialogTitle>
        </AlertDialogHeader>
        <AlertDialogDescription>
          <div className='flex flex-col w-full'>
            <div className='flex justify-between'>
              <label>Trace Number:</label>
              <span>{atmResponse?.traceNumber}</span>
            </div>
            <div className='flex justify-between'>
              <label>Response Code:</label>
              <span>{atmResponse?.responseCode}</span>
            </div>
            <div className='flex justify-between'>
              <label>RRN:</label>
              <span>{atmResponse?.rrn}</span>
            </div>
            <div className='flex justify-between'>
              <label>Balance:</label>
              <span>{atmResponse?.balance}</span>
            </div>
          </div>
        </AlertDialogDescription>
        <AlertDialogFooter>
          <AlertDialogCancel>Close</AlertDialogCancel>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default AtmResponseDialog