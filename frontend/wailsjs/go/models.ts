export namespace internal {
	
	export class Bird {
	    speciesCode: string;
	    sciName: string;
	    comName: string;
	
	    static createFrom(source: any = {}) {
	        return new Bird(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.speciesCode = source["speciesCode"];
	        this.sciName = source["sciName"];
	        this.comName = source["comName"];
	    }
	}

}

export namespace maps {
	
	export class Observation {
	    speciesCode: string;
	    comName: string;
	    sciName: string;
	    locId: string;
	    locName: string;
	    obsDt: string;
	    lat: number;
	    lng: number;
	
	    static createFrom(source: any = {}) {
	        return new Observation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.speciesCode = source["speciesCode"];
	        this.comName = source["comName"];
	        this.sciName = source["sciName"];
	        this.locId = source["locId"];
	        this.locName = source["locName"];
	        this.obsDt = source["obsDt"];
	        this.lat = source["lat"];
	        this.lng = source["lng"];
	    }
	}
	export class Location {
	    locId: string;
	    locName: string;
	    lat: number;
	    lng: number;
	    observations: Observation[];
	
	    static createFrom(source: any = {}) {
	        return new Location(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.locId = source["locId"];
	        this.locName = source["locName"];
	        this.lat = source["lat"];
	        this.lng = source["lng"];
	        this.observations = this.convertValues(source["observations"], Observation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

