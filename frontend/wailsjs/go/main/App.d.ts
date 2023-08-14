// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function CloseTunnel():Promise<void>;

export function GetConfigs():Promise<Array<main.Config>>;

export function GetMessages(arg1:number):Promise<Array<main.Message>>;

export function PingTunnel():Promise<void>;

export function SendMessage(arg1:main.Message):Promise<main.AtmResponse>;

export function SendReversalMessage(arg1:number):Promise<main.AtmResponse>;

export function UpdateConfigs(arg1:Array<main.Config>):Promise<void>;

export function UseTunnel():Promise<void>;
