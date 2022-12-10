export namespace backend {
	
	export class DeviceInformationJson {
	    path: string;
	    vendor_id: number;
	    product_id: number;
	    serial: string;
	    manufacturer: string;
	    product: string;
	    interface: number;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInformationJson(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.vendor_id = source["vendor_id"];
	        this.product_id = source["product_id"];
	        this.serial = source["serial"];
	        this.manufacturer = source["manufacturer"];
	        this.product = source["product"];
	        this.interface = source["interface"];
	    }
	}

}

