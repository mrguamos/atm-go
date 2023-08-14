// export type Message = {
//     transaction: Transaction
//     switch: AtmSwitch
//     primaryAccountNumber: string
//     transactionAmount: number
//     acquiringInstitutionCode: string
//     receivingInstitutionCode: string
//     transactionFee: number
//     terminalNameAndLocation: string
//     currencyCode: Currency
//     terminalId: string
//     sourceAccount: string
//     destinationAccount: string
//     channel: Channel
//     device: Device
//     targetBank: Bank
//     rrn: string    
// 	traceNumber: string    
// 	transmissionDateTime: string    
// 	localTransactionDateTime: string    
// 	originalDataElements: string     
// 	mti: string   
// }

// export type AtmResponse = {
//     traceNumber: string
//     responseCode: string
//     balance?: string
//     rrn?: string
//     err?: string
// }

export enum Currency {
    PHP = '608',
    USD = '840'
}

export enum Device {
    ATM = '6011',
    POS = '6012',
    NAD = '6016'
}

export enum AtmSwitch {
    CORTEX = 'CORTEX',
    POSTBRIDGE = 'POSTBRIDGE',
    COREWARE = 'COREWARE'
}

export enum Transaction {
    WITHDRAW = 'WITHDRAW',
    BAL_INQ = 'BAL_INQ',
    FT = 'FT',
    IBFTC = 'IBFTC',
    IBFTD = 'IBFTD',
    ELOAD = 'ELOAD',
    BILLS = 'BILLS',
    PURCHASE = 'PURCHASE'
}

export enum Channel {
    ON_US = 'ON_US',
    OFF_US = 'OFF_US',
    MASTERCARD = 'MASTERCARD'
}

export enum Bank {
    OTHER_BANK = 'OTHER_BANK',
    INTER_SYSTEM = 'INTER_SYSTEM',
}

export enum MTI {
    FINANCIAL_REQUEST_MASTER_VISA = 100,
    FINANCIAL_REQUEST = 200,
    FINANCIAL_ADVICE = 220,
    FINANCIAL_REVERSAL = 400,
    FINANCIAL_REVERSAL_ADVICE = 420,
    FINANCIAL_REVERSAL_REPEAT_ADVICE = 421,
    NETWORK_MANAGEMENT_REQUEST = 800

}

