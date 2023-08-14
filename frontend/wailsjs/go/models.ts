export namespace main {
	
	export class AtmResponse {
	    traceNumber: string;
	    responseCode: string;
	    balance: string;
	    rrn: string;
	
	    static createFrom(source: any = {}) {
	        return new AtmResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.traceNumber = source["traceNumber"];
	        this.responseCode = source["responseCode"];
	        this.balance = source["balance"];
	        this.rrn = source["rrn"];
	    }
	}
	export class Config {
	    key: string;
	    value?: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}
	export class Message {
	    transaction: string;
	    switch: string;
	    primaryAccountNumber?: string;
	    transactionAmount?: number;
	    acquiringInstitutionCode: string;
	    receivingInstitutionCode?: string;
	    transactionFee?: number;
	    terminalNameAndLocation: string;
	    currencyCode: string;
	    terminalId: string;
	    sourceAccount?: string;
	    destinationAccount?: string;
	    channel: string;
	    device: string;
	    targetBank?: string;
	    id?: number;
	    rrn?: string;
	    traceNumber?: string;
	    transmissionDateTime?: string;
	    localTransactionDateTime?: string;
	    originalDataElements?: string;
	    mti?: string;
	    processCod?: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.transaction = source["transaction"];
	        this.switch = source["switch"];
	        this.primaryAccountNumber = source["primaryAccountNumber"];
	        this.transactionAmount = source["transactionAmount"];
	        this.acquiringInstitutionCode = source["acquiringInstitutionCode"];
	        this.receivingInstitutionCode = source["receivingInstitutionCode"];
	        this.transactionFee = source["transactionFee"];
	        this.terminalNameAndLocation = source["terminalNameAndLocation"];
	        this.currencyCode = source["currencyCode"];
	        this.terminalId = source["terminalId"];
	        this.sourceAccount = source["sourceAccount"];
	        this.destinationAccount = source["destinationAccount"];
	        this.channel = source["channel"];
	        this.device = source["device"];
	        this.targetBank = source["targetBank"];
	        this.id = source["id"];
	        this.rrn = source["rrn"];
	        this.traceNumber = source["traceNumber"];
	        this.transmissionDateTime = source["transmissionDateTime"];
	        this.localTransactionDateTime = source["localTransactionDateTime"];
	        this.originalDataElements = source["originalDataElements"];
	        this.mti = source["mti"];
	        this.processCod = source["processCod"];
	    }
	}

}

