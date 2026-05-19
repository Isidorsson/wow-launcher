export namespace config {
	
	export class Server {
	    id: string;
	    name: string;
	    manifestUrl: string;
	    website: string;
	    newsFeedUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new Server(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.manifestUrl = source["manifestUrl"];
	        this.website = source["website"];
	        this.newsFeedUrl = source["newsFeedUrl"];
	    }
	}

}

export namespace main {
	
	export class BrandingDTO {
	    launcherName: string;
	    windowTitle: string;
	    primaryColor: string;
	
	    static createFrom(source: any = {}) {
	        return new BrandingDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.launcherName = source["launcherName"];
	        this.windowTitle = source["windowTitle"];
	        this.primaryColor = source["primaryColor"];
	    }
	}
	export class DetectedInstall {
	    root: string;
	    locale: string;
	
	    static createFrom(source: any = {}) {
	        return new DetectedInstall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.root = source["root"];
	        this.locale = source["locale"];
	    }
	}
	export class ProfileDTO {
	    serverId: string;
	    root: string;
	    locale: string;
	    exists: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProfileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serverId = source["serverId"];
	        this.root = source["root"];
	        this.locale = source["locale"];
	        this.exists = source["exists"];
	    }
	}

}

