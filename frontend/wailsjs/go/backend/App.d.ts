// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {context} from '../models';
import {backend} from '../models';

export function BeforeClose(arg1:context.Context):Promise<boolean>;

export function DomReady(arg1:context.Context):Promise<void>;

export function GetDeviceInformation():Promise<void>;

export function LoadAvaliableDevices():Promise<Array<backend.DeviceInformationJson>>;

export function LoadDevices(arg1:context.Context):Promise<void>;

export function LoadDevicesPolling():Promise<Error>;

export function SelectDevice(arg1:string):Promise<Error>;

export function SetButtonWrapper(arg1:string,arg2:string):Promise<void>;

export function SetDpiWrapper(arg1:number):Promise<void>;

export function SetLEDWrapper(arg1:string):Promise<void>;

export function Shutdown(arg1:context.Context):Promise<void>;

export function UpdateAvaliableDevices(arg1:Array<backend.DeviceInformationJson>):Promise<void>;
