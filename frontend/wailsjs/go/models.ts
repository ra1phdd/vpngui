export namespace models {
	
	export class Config {
	    ActiveVPN: boolean;
	    DisableRoutes: boolean;
	    ListMode: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ActiveVPN = source["ActiveVPN"];
	        this.DisableRoutes = source["DisableRoutes"];
	        this.ListMode = source["ListMode"];
	    }
	}
	export class InboundConfig {
	    listen: string;
	    port: number;
	    protocol: string;
	    settings?: {[key: string]: any};
	    streamSettings?: {[key: string]: any};
	    tag?: string;
	
	    static createFrom(source: any = {}) {
	        return new InboundConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.listen = source["listen"];
	        this.port = source["port"];
	        this.protocol = source["protocol"];
	        this.settings = source["settings"];
	        this.streamSettings = source["streamSettings"];
	        this.tag = source["tag"];
	    }
	}
	export class Rule {
	    RuleType: string;
	    RuleValue: string;
	
	    static createFrom(source: any = {}) {
	        return new Rule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RuleType = source["RuleType"];
	        this.RuleValue = source["RuleValue"];
	    }
	}
	export class ListConfig {
	    Type: string;
	    Rules: Rule[];
	    DomainStrategy: string;
	    DomainMatcher: string;
	
	    static createFrom(source: any = {}) {
	        return new ListConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Rules = this.convertValues(source["Rules"], Rule);
	        this.DomainStrategy = source["DomainStrategy"];
	        this.DomainMatcher = source["DomainMatcher"];
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
	export class LogConfig {
	    loglevel: string;
	
	    static createFrom(source: any = {}) {
	        return new LogConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.loglevel = source["loglevel"];
	    }
	}
	export class OutboundConfig {
	    sendThrough?: string;
	    protocol: string;
	    settings?: {[key: string]: any};
	    tag: string;
	    streamSettings?: {[key: string]: any};
	    mux?: any;
	
	    static createFrom(source: any = {}) {
	        return new OutboundConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sendThrough = source["sendThrough"];
	        this.protocol = source["protocol"];
	        this.settings = source["settings"];
	        this.tag = source["tag"];
	        this.streamSettings = source["streamSettings"];
	        this.mux = source["mux"];
	    }
	}
	export class RoutingRule {
	    type: string;
	    domain?: string[];
	    ip?: string[];
	    port?: string;
	    outboundTag: string;
	
	    static createFrom(source: any = {}) {
	        return new RoutingRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.domain = source["domain"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.outboundTag = source["outboundTag"];
	    }
	}
	export class RoutingConfig {
	    domainStrategy: string;
	    domainMatcher: string;
	    rules?: RoutingRule[];
	    balancers?: string[];
	
	    static createFrom(source: any = {}) {
	        return new RoutingConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domainStrategy = source["domainStrategy"];
	        this.domainMatcher = source["domainMatcher"];
	        this.rules = this.convertValues(source["rules"], RoutingRule);
	        this.balancers = source["balancers"];
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
	
	
	export class Xray {
	    log: LogConfig;
	    inbounds: InboundConfig[];
	    outbounds: OutboundConfig[];
	    routing?: RoutingConfig;
	    stats: any;
	
	    static createFrom(source: any = {}) {
	        return new Xray(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.log = this.convertValues(source["log"], LogConfig);
	        this.inbounds = this.convertValues(source["inbounds"], InboundConfig);
	        this.outbounds = this.convertValues(source["outbounds"], OutboundConfig);
	        this.routing = this.convertValues(source["routing"], RoutingConfig);
	        this.stats = source["stats"];
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

