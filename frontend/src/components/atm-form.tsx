'use client'

import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import * as z from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { AtmSwitch, Bank, Channel, Currency, Device, Transaction } from '@/lib/message'
import { useForm } from 'react-hook-form'
import { enumFromStringValue, getEnumKeys } from '@/lib/helper'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from './ui/form'
import { useState } from 'react'
import { SendFinancialMessage } from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'
import AtmResponseDialog from './atm-response-dialog'
import { useRecoilState } from 'recoil'
import { keyState, loadingState, messageState } from '@/store/state'
import { useToast } from './ui/use-toast'


export function AtmForm () {
  const formSchema = z.object({
    transaction: z.nativeEnum(Transaction, {
      invalid_type_error: `Invalid transaction type, should be any of ${getEnumKeys(Transaction)}.`,
      required_error: 'Transaction is required.'
    }),
    switch: z.nativeEnum(AtmSwitch, {
      invalid_type_error: `Invalid transaction type should be any of ${getEnumKeys(AtmSwitch)}`,
      required_error: 'Switch is required.'
    }),
    device: z.nativeEnum(Device, {
      invalid_type_error: `Invalid device, should be any of ${getEnumKeys(Device)}.`,
      required_error: 'Device is required.'
    }),
    primaryAccountNumber: z.string({}).max(19).optional(),
    transactionAmount: z.coerce.number().optional(),
    acquiringInstitutionCode: z.string({}).min(1, 'Acquiring Institution Code is required.').max(11),
    receivingInstitutionCode: z.string({}).max(11).optional(),
    transactionFee: z.coerce.number().optional(),
    terminalNameAndLocation: z.string({}).max(99),
    currencyCode: z.nativeEnum(Currency, {
      invalid_type_error: `Invalid currency, should be any of ${getEnumKeys(Currency)}.`,
      required_error: 'Currency is required.'
    }),
    terminalId: z.string({}).min(8, 'Terminal ID must contain 8 characters.').max(8),
    sourceAccount: z.string().max(28).optional(),
    destinationAccount: z.string().max(28).optional(),
    channel: z.nativeEnum(Channel, {
      invalid_type_error: `Invalid channel, should be any of ${getEnumKeys(Channel)}.`,
      required_error: 'Channel is required.'
    }),
    targetBank: z.nativeEnum(Bank, {
      invalid_type_error: `Invalid target bank, should be any of ${getEnumKeys(Bank)}.`
    }).optional()
  })

  const [message, setMessage] = useRecoilState(messageState)
  const [, setKey] = useRecoilState(keyState)

  const reset = () => {
    setMessage({} as main.Message)
    setKey(new Date().getMilliseconds())
  }

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      transaction: enumFromStringValue(Transaction, message.transaction),
      switch: enumFromStringValue(AtmSwitch, message.switch),
      device: enumFromStringValue(Device, message.device),
      channel: enumFromStringValue(Channel, message.channel),
      primaryAccountNumber: message.primaryAccountNumber ?? '',
      transactionAmount: message.transactionAmount ?? 0,
      transactionFee: message.transactionFee ?? 0,
      acquiringInstitutionCode: message.acquiringInstitutionCode ?? '',
      receivingInstitutionCode: message.receivingInstitutionCode ?? '',
      currencyCode: enumFromStringValue(Currency, message.currencyCode),
      destinationAccount: message.destinationAccount ?? '',
      sourceAccount: message.sourceAccount ?? '',
      targetBank: enumFromStringValue(Bank, message.targetBank ?? ''),
      terminalId:  message.terminalId ?? '',
      terminalNameAndLocation: message.terminalNameAndLocation ?? '',
    }
  })

  const [, setLoading] = useRecoilState(loadingState)

  const [atmResponse, setAtmResponse] = useState<main.AtmResponse>()
  const [isOpen, setOpen] = useState(false)
  const { toast } = useToast()
  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    let response
    try {
      setLoading(true)
      response = await SendFinancialMessage(data)
      setAtmResponse(response)
      setOpen(true)
    }catch(error: any) {
      toast({
        description: error,
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <AtmResponseDialog isOpen={isOpen} setOpen={setOpen} atmResponse={atmResponse}/>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
          <Card>
            <CardHeader>
              <CardTitle className="text-center">ATM Terminal</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid w-full items-start gap-x-10 gap-y-4 sm:grid-cols-2">
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="transaction"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="transaction">Transaction</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <SelectTrigger id="transaction" aria-controls='transaction'>
                              <SelectValue placeholder="Select Transaction"/>
                            </SelectTrigger>
                            <SelectContent position="popper">
                              <SelectItem value="WITHDRAW">Cash Withdrawal</SelectItem>
                              <SelectItem value="BAL_INQ">Balance Inquiry</SelectItem>
                              <SelectItem value="FT">Fund Transfer</SelectItem>
                              <SelectItem value="IBFTC">IBFT Credit</SelectItem>
                              <SelectItem value="IBFTD">IBFT Debit</SelectItem>
                              <SelectItem value="ELOAD">E-Load</SelectItem>
                              <SelectItem value="BILLS">Bills Payment</SelectItem>
                              <SelectItem value="PURCHASE">Purchase</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="switch"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="switch">Switch</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <SelectTrigger id="switch" aria-controls='switch'>
                              <SelectValue placeholder="Select Switch"/>
                            </SelectTrigger>
                            <SelectContent position="popper">
                              <SelectItem value="CORTEX">Cortex</SelectItem>
                              <SelectItem value="POSTBRIDGE">Postbridge</SelectItem>
                              <SelectItem value="COREWARE">Coreware</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="device"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="device">Device</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <SelectTrigger id="device" aria-controls='device'>
                              <SelectValue placeholder="Select Device" />
                            </SelectTrigger>
                            <SelectContent position="popper">
                              <SelectItem value="6011">ATM</SelectItem>
                              <SelectItem value="6012">POS</SelectItem>
                              <SelectItem value="6016">NAD</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="primaryAccountNumber"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="pan">Primary Account Number</FormLabel>
                        <FormControl>
                          <Input id="pan" aria-describedby="pan" max={19} placeholder="5395998824355777" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="transactionAmount"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="amount">Transaction Amount</FormLabel>
                        <FormControl>
                          <Input aria-describedby="amount" id="amount" min="0" type="number" step="0.01" placeholder="100.00" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="transactionFee"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="fee">Transaction Fee</FormLabel>
                        <FormControl>
                          <Input aria-describedby="fee" id="fee" min="0" type="number" step="0.01" placeholder="20.00" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="currencyCode"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="currency">Currency</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <SelectTrigger id="currency" aria-controls='currency'>
                              <SelectValue placeholder="Select Currency" />
                            </SelectTrigger>
                            <SelectContent position="popper">
                              <SelectItem value="608">PHP</SelectItem>
                              <SelectItem value="840">USD</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="channel"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="channel">Channel</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <SelectTrigger id="channel" aria-controls='channel'>
                              <SelectValue placeholder="Select Channel" />
                            </SelectTrigger>
                            <SelectContent position="popper">
                              <SelectItem value="ON_US">On Us</SelectItem>
                              <SelectItem value="OFF_US">Off Us</SelectItem>
                              <SelectItem value="MASTERCARD">Mastercard</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="acquiringInstitutionCode"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="aiic">Acquiring Institution Code</FormLabel>
                        <FormControl>
                          <Input id="aiic" aria-describedby="aiic" max={11} placeholder="928" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="receivingInstitutionCode"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="riic">Receiving Institution Code</FormLabel>
                        <FormControl>
                          <Input id="riic" aria-describedby="riic" max={11} placeholder="908" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="terminalId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="tid">Terminal ID</FormLabel>
                        <FormControl>
                          <Input id="tid" aria-describedby="tid" max={8} placeholder="61740007" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="terminalNameAndLocation"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="tname">Terminal Name and Location</FormLabel>
                        <FormControl>
                          <Input id="tname" aria-describedby="tname" max={99} placeholder="BGC ATM1 TAGUIG PH" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="sourceAccount"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="source">Source Account</FormLabel>
                        <FormControl>
                          <Input id="source" aria-describedby="source" max={28} placeholder="000100000001366" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="destinationAccount"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="destination">Destination Account</FormLabel>
                        <FormControl>
                          <Input id="destination" aria-describedby="destination" max={28} placeholder="000100000001366" {...field}/>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="targetBank"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel htmlFor="target">Target Bank</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <SelectTrigger id="target" aria-controls='target'>
                              <SelectValue placeholder="Select Target Bank"/>
                            </SelectTrigger>
                            <SelectContent position="popper">
                              <SelectItem value="OTHER_BANK">Other Bank</SelectItem>
                              <SelectItem value="INTER_SYSTEM">Inter System</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>
            </CardContent>
            <hr className='my-10 mx-10 h-1 rounded bg-slate-700' />
            <CardFooter className="flex justify-between">
              <Button variant={'destructive'} type='button' onClick={reset}>Clear</Button>
              <Button type='submit'>Submit</Button>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </>
  )
}
